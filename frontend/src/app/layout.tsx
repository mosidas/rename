import type { Metadata } from 'next';
import Script from 'next/script';
import { Noto_Sans_JP } from 'next/font/google';
import Header from '../components/Header';
import './globals.css';

const notoSansJp = Noto_Sans_JP({
  variable: '--font-sans',
  display: 'swap',
  subsets: ['latin'],
  weight: ['400', '700'],
});

export const metadata: Metadata = {
  title: 'File Rename - macOS専用ファイル一括リネームツール',
  description: 'macOS専用ファイル一括リネームツール。正規表現対応、履歴機能、リアルタイムプレビューを搭載。',
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja" suppressHydrationWarning>
      <head>
        <Script id="theme-init" strategy="beforeInteractive">
          {`(function(){
            try {
              var key='theme';
              var stored=localStorage.getItem(key);
              if (stored === 'light' || stored === 'dark') {
                document.documentElement.setAttribute('data-theme', stored);
              } else {
                // system: 属性を外して OS の設定に追従
                document.documentElement.removeAttribute('data-theme');
              }
            } catch (e) {}
          })();`}
        </Script>
      </head>
      <body className={`${notoSansJp.variable} antialiased`}>
        <Header />
        {children}
      </body>
    </html>
  );
}
