"use client";

import Image from "next/image";
import { useEffect, useRef, useState } from "react";

type Theme = "light" | "dark" | "system";
const STORAGE_KEY = "theme";

function getCurrentTheme(): Theme {
  if (typeof document === "undefined") return "system";
  const attr = document.documentElement.getAttribute("data-theme");
  if (attr === "light" || attr === "dark") return attr;
  return "system";
}

function applyTheme(theme: Theme) {
  if (theme === "light" || theme === "dark") {
    document.documentElement.setAttribute("data-theme", theme);
    try {
      localStorage.setItem(STORAGE_KEY, theme);
    } catch {}
  } else {
    document.documentElement.removeAttribute("data-theme");
    try {
      localStorage.setItem(STORAGE_KEY, "system");
    } catch {}
  }
}

export default function Header() {
  const [open, setOpen] = useState(false);
  const [theme, setTheme] = useState<Theme>("system");
  const menuRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    setTheme(getCurrentTheme());
  }, []);

  useEffect(() => {
    const onClickAway = (e: MouseEvent) => {
      if (!open) return;
      const target = e.target as Node;
      if (menuRef.current && !menuRef.current.contains(target)) {
        setOpen(false);
      }
    };
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") setOpen(false);
    };
    document.addEventListener("click", onClickAway);
    document.addEventListener("keydown", onKey);
    return () => {
      document.removeEventListener("click", onClickAway);
      document.removeEventListener("keydown", onKey);
    };
  }, [open]);

  const choose = (t: Theme) => {
    setTheme(t);
    applyTheme(t);
    setOpen(false);
  };

  return (
    <header className="sticky top-0 z-40 border-b bg-background/80 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="mx-auto max-w-5xl px-4 h-12 flex items-center justify-between">
        <div className="text-sm font-medium">App</div>
        <div className="relative" ref={menuRef}>
          <button
            type="button"
            aria-haspopup="menu"
            aria-expanded={open}
            aria-label="設定"
            className="p-2 rounded-md border hover:bg-foreground/5"
            onClick={() => setOpen((v) => !v)}
          >
            <Image
              src="/vscode-icons--file-type-config.svg"
              alt="設定"
              width={18}
              height={18}
              priority
            />
          </button>

          {open && (
            <div
              role="menu"
              className="absolute right-0 mt-2 w-40 rounded-md border bg-background shadow-md p-1 text-sm"
            >
              {[
                { key: "light", label: "ライト" },
                { key: "dark", label: "ダーク" },
                { key: "system", label: "システム" },
              ].map((opt) => (
                <button
                  key={opt.key}
                  role="menuitemradio"
                  aria-checked={theme === (opt.key as Theme)}
                  className={`w-full text-left px-3 py-2 rounded hover:bg-foreground/5 ${
                    theme === (opt.key as Theme) ? "bg-foreground/10" : ""
                  }`}
                  onClick={() => choose(opt.key as Theme)}
                >
                  {opt.label}
                </button>
              ))}
            </div>
          )}
        </div>
      </div>
    </header>
  );
}

