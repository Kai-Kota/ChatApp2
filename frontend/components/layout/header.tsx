"use client";

import { useEffect, useState } from "react";

export default function Header() {
  const [userName, setUserName] = useState<string>("");

  useEffect(() => {
    if (typeof window !== "undefined") {
      const name = localStorage.getItem("userName") || "";
      setUserName(name);
    }
  }, []);

  const displayName = userName || "未ログイン";

  return (
    <header className="w-full h-16 bg-cyan-800 shadow-md text-white flex items-center px-4">
      <h1 className="text-xl font-bold">Chatty</h1>
      <div className="ml-auto flex items-center gap-3">
        <span className="rounded-full bg-cyan-600 w-8 h-8" aria-hidden="true" />
        <span className="text-sm">{displayName}</span>
      </div>
    </header>
  );
}