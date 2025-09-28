'use client';

import Image from 'next/image';
import { useState } from 'react';
import { Greet } from '../../wailsjs/go/main/App';

export default function Home() {
  return (
    <div className="min-h-screen p-6">
      <main className="flex flex-col items-center gap-6 w-full max-w-md mx-auto pt-12">
        <div>
          あのイーハトーヴォのすきとおった風、夏でも底に冷たさをもつ青いそら、うつくしい森で飾られたモリーオ市、郊外のぎらぎらひかる草の波。
        </div>
        <div>
          Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
        </div>
      </main>
    </div>
  );
}
