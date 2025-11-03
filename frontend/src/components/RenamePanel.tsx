'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import { SelectFiles, GeneratePreview, ExecuteRename, GetHistory, GetInitialFiles } from '../../wailsjs/go/main/App';
import { EventsOn } from '../../wailsjs/runtime/runtime';
import { main, domain } from '../../wailsjs/go/models';

// Use Wails-generated types
type FilePreview = main.FilePreview;
type HistoryEntry = domain.HistoryEntry;

// Constants
const PREVIEW_DEBOUNCE_MS = 300;
const MAX_HISTORY_DISPLAY = 10;

export default function RenamePanel() {
  const [selectedFiles, setSelectedFiles] = useState<string[]>([]);
  const [pattern, setPattern] = useState('');
  const [replacement, setReplacement] = useState('');
  const [isRegex, setIsRegex] = useState(false);
  const [caseInsensitive, setCaseInsensitive] = useState(false);
  const [previews, setPreviews] = useState<FilePreview[]>([]);
  const [history, setHistory] = useState<HistoryEntry[]>([]);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');
  const [showHistoryDropdown, setShowHistoryDropdown] = useState(false);
  const patternInputRef = useRef<HTMLInputElement>(null);

  // Load history and initial files on mount
  useEffect(() => {
    loadHistory();

    // Check if files were provided on startup via command-line
    GetInitialFiles().then((files) => {
      if (files && files.length > 0) {
        setSelectedFiles(files);
        setMessage(`${files.length}個のファイルを選択しました`);
      }
    });

    // Listen for files loaded from second instance (when app is already running)
    const cleanup = EventsOn('files:loaded', (files: string[]) => {
      if (files && files.length > 0) {
        setSelectedFiles(files);
        setMessage(`${files.length}個のファイルを選択しました`);
      }
    });

    // Cleanup event listener on unmount
    return cleanup;
  }, []);

  // Auto-generate preview when inputs change
  useEffect(() => {
    if (selectedFiles.length > 0) {
      generatePreviewDebounced();
    }
  }, [pattern, replacement, isRegex, caseInsensitive, selectedFiles]);

  const loadHistory = async () => {
    try {
      const entries = await GetHistory();
      setHistory(entries || []);
    } catch (err) {
      // Ignore history load errors
    }
  };

  const handleSelectFiles = async () => {
    try {
      const files = await SelectFiles();
      if (files && files.length > 0) {
        setSelectedFiles(files);
        setMessage(`${files.length}個のファイルを選択しました`);
      }
    } catch (err) {
      setMessage('ファイル選択に失敗しました');
    }
  };

  // Debounced preview generation
  const generatePreviewDebounced = useCallback((() => {
    let timeoutId: NodeJS.Timeout;
    return () => {
      clearTimeout(timeoutId);
      timeoutId = setTimeout(async () => {
        if (selectedFiles.length === 0) {
          return;
        }

        try {
          const result = await GeneratePreview(pattern, replacement, isRegex, caseInsensitive);
          setPreviews(result || []);
          const changedCount = result?.filter(p => p.hasChanged).length || 0;
          if (changedCount > 0) {
            setMessage(`${changedCount}個のファイルが変更されます`);
          }
        } catch (err: any) {
          setMessage(`プレビュー生成エラー: ${err.message || '不明なエラー'}`);
          setPreviews([]);
        }
      }, PREVIEW_DEBOUNCE_MS);
    };
  })(), [selectedFiles, pattern, replacement, isRegex, caseInsensitive]);

  const handleExecuteRename = async () => {
    if (previews.length === 0) {
      setMessage('プレビューを生成してください');
      return;
    }

    const changedCount = previews.filter(p => p.hasChanged).length;
    if (changedCount === 0) {
      setMessage('変更するファイルがありません');
      return;
    }

    setLoading(true);
    try {
      const result = await ExecuteRename();

      // Update selectedFiles with new paths after rename
      if (result.NewFilePaths && result.NewFilePaths.length > 0) {
        setSelectedFiles(result.NewFilePaths);
      }

      // Reload history (backend adds to history automatically)
      if (result.SuccessCount > 0) {
        await loadHistory();
      }

      setMessage(
        `成功: ${result.SuccessCount}件, 失敗: ${result.FailureCount}件` +
        (result.Errors && result.Errors.length > 0 ? `\nエラー: ${result.Errors.join(', ')}` : '')
      );

      // Don't reset - keep inputs and files for continuous renaming
    } catch (err: any) {
      setMessage(`リネーム実行エラー: ${err.message || '不明なエラー'}`);
    } finally {
      setLoading(false);
    }
  };

  const handleHistorySelect = (entry: HistoryEntry) => {
    setPattern(entry.pattern);
    setReplacement(entry.replacement);
    setIsRegex(entry.isRegex);
    setCaseInsensitive(entry.caseInsensitive);
    setShowHistoryDropdown(false);
  };

  return (
    <div className="h-[calc(100vh-3rem)] flex flex-col p-6 bg-muted overflow-hidden">
      {/* Header */}
      <div className="mb-6 flex items-center gap-4">
        <button
          onClick={handleSelectFiles}
          className="px-4 py-2 bg-accent text-accent-foreground rounded hover:bg-accent/90 transition"
          disabled={loading}
        >
          ファイルを選択
        </button>
        {message && (
          <span className="text-muted-foreground text-sm">
            {message}
          </span>
        )}
      </div>

      {/* 2 Column Layout - 1:2 ratio */}
      <div className="flex-1 flex gap-6 overflow-hidden">
        {/* Left Column - Input (flex-1 for 1 part) */}
        <div className="flex-1 bg-background p-6 rounded-lg border shadow overflow-auto">
          <div className="space-y-4">
            {/* Pattern Input with History Dropdown */}
            <div className="relative">
              <label className="block text-sm font-medium mb-2 text-foreground">
                置換前
              </label>
              <input
                ref={patternInputRef}
                type="text"
                value={pattern}
                onChange={(e) => setPattern(e.target.value)}
                onFocus={() => setShowHistoryDropdown(true)}
                onBlur={() => setTimeout(() => setShowHistoryDropdown(false), 200)}
                className="w-full px-3 py-2 border rounded bg-background text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-accent"
                placeholder="例: test"
              />

              {/* History Dropdown */}
              {showHistoryDropdown && history.length > 0 && (
                <div className="absolute z-10 w-full mt-1 bg-background border rounded shadow-lg max-h-60 overflow-auto">
                  {history.slice(0, MAX_HISTORY_DISPLAY).map((entry, index) => (
                    <button
                      key={index}
                      onClick={() => handleHistorySelect(entry)}
                      className="w-full text-left px-3 py-2 hover:bg-muted border-b last:border-b-0"
                    >
                      <div className="flex flex-wrap items-center gap-2 text-sm">
                        <span className="font-mono text-foreground">{entry.pattern}</span>
                        <span className="text-xs bg-muted-foreground/20 text-muted-foreground px-2 py-0.5 rounded">→</span>
                        <span className="font-mono text-foreground">{entry.replacement}</span>
                      </div>
                      <div className="flex gap-2 mt-1">
                        {entry.isRegex && (
                          <span className="text-xs bg-accent/20 text-accent px-2 py-0.5 rounded">
                            正規表現
                          </span>
                        )}
                        {entry.caseInsensitive && (
                          <span className="text-xs bg-success/20 text-success px-2 py-0.5 rounded">
                            大小文字無視
                          </span>
                        )}
                      </div>
                    </button>
                  ))}
                </div>
              )}
            </div>

            {/* Replacement Input */}
            <div>
              <label className="block text-sm font-medium mb-2 text-foreground">
                置換後
              </label>
              <input
                type="text"
                value={replacement}
                onChange={(e) => setReplacement(e.target.value)}
                className="w-full px-3 py-2 border rounded bg-background text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-accent"
                placeholder="例: TEST"
              />
            </div>

            {/* Options */}
            <div className="space-y-2">
              <label className="flex items-center cursor-pointer">
                <input
                  type="checkbox"
                  checked={isRegex}
                  onChange={(e) => setIsRegex(e.target.checked)}
                  className="mr-2 w-4 h-4 rounded border accent-checkbox"
                />
                <span className="text-sm text-foreground">正規表現</span>
              </label>
              <label className="flex items-center cursor-pointer">
                <input
                  type="checkbox"
                  checked={caseInsensitive}
                  onChange={(e) => setCaseInsensitive(e.target.checked)}
                  className="mr-2 w-4 h-4 rounded border accent-checkbox"
                />
                <span className="text-sm text-foreground">
                  大文字小文字を区別しない
                </span>
              </label>
            </div>

            {/* Execute Button */}
            <button
              onClick={handleExecuteRename}
              disabled={loading || previews.length === 0 || previews.filter(p => p.hasChanged).length === 0}
              className="w-full px-4 py-3 bg-destructive text-accent-foreground rounded hover:bg-destructive/90 transition disabled:opacity-50 disabled:cursor-not-allowed font-medium"
            >
              リネーム実行
            </button>
          </div>
        </div>

        {/* Right Column - Preview (flex-[2] for 2 parts) */}
        <div className="flex-[2] bg-background rounded-lg border shadow overflow-hidden flex flex-col">
          <div className="flex-1 overflow-auto">
            <table className="w-full">
              <thead className="bg-muted sticky top-0">
                <tr>
                  <th className="px-4 py-2 text-left text-sm font-medium text-foreground">
                    元のファイル名
                  </th>
                  <th className="px-4 py-2 text-left text-sm font-medium text-foreground">
                    新しいファイル名
                  </th>
                  <th className="px-4 py-2 text-center text-sm font-medium text-foreground w-16">
                    変更
                  </th>
                </tr>
              </thead>
              <tbody>
                {selectedFiles.length === 0 ? (
                  <tr>
                    <td colSpan={3} className="px-4 py-8 text-center text-muted-foreground">
                      ファイルを選択してください
                    </td>
                  </tr>
                ) : (
                  previews.map((preview, index) => (
                    <tr
                      key={index}
                      className={`border-t ${
                        preview.hasChanged ? 'bg-accent/5' : ''
                      }`}
                    >
                      <td className="px-4 py-2 text-sm text-foreground">
                        {preview.originalName}
                      </td>
                      <td className="px-4 py-2 text-sm text-foreground font-medium">
                        {preview.newName}
                      </td>
                      <td className="px-4 py-2 text-sm text-center">
                        {preview.hasChanged ? (
                          <span className="text-success">✓</span>
                        ) : (
                          <span className="text-muted-foreground">-</span>
                        )}
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
}
