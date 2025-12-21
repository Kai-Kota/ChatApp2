"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Header from "@/components/layout/header";

export default function SignupPage() {
    const router = useRouter();
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        setSuccess(null);

        if (username.trim().length < 3) {
            setError("ユーザー名は3文字以上にしてください。");
            return;
        }
        if (password.length < 6) {
            setError("パスワードは6文字以上にしてください。");
            return;
        }

        setLoading(true);
        try {
            const res = await fetch("http://localhost:8080/login", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ username, password }),
                credentials: "include",
            });

            if (!res.ok) {
                const data = await res.json().catch(() => ({}));
                setError(data.message || `エラー: ${res.status}`);
                setLoading(false);
                return;
            }
            console.log(await(res).json())
            if (res.ok) console.log("Signup successful");

            setSuccess("ログイン中...");
            setUsername("");
            setPassword("");
            setTimeout(() => router.push("/home"), 1200);
        } catch (err) {
            setError("ネットワークエラーが発生しました。");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen flex flex-col bg-gray-50">
            <Header />
            <div className="flex-auto flex items-center justify-center">
                <div className="w-full max-w-md bg-white rounded-lg shadow-md p-6">
                    <h1 className="text-2xl font-semibold mb-4">ログイン</h1>

                    <form onSubmit={handleSubmit} className="space-y-4">
                        <div>
                            <label className="block text-sm font-medium mb-1">ユーザー名</label>
                            <input
                                type="text"
                                value={username}
                                onChange={(e) => setUsername(e.target.value)}
                                className="w-full rounded-md border px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-300"
                                placeholder="ユーザー名"
                                required
                            />
                        </div>
        
                        <div>
                            <label className="block text-sm font-medium mb-1">パスワード</label>
                            <input
                                type="password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                className="w-full rounded-md border px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-300"
                                placeholder="パスワード"
                                required
                            />
                        </div>
        
                        {error && <div className="text-sm text-red-600">{error}</div>}
                        {success && <div className="text-sm text-green-600">{success}</div>}
        
                        <div className="flex items-center justify-between">
                            <button
                                type="submit"
                                className="bg-blue-500 text-white px-4 py-2 rounded-md hover:bg-blue-600 disabled:opacity-60"
                                disabled={loading}
                            >
                                {loading ? "送信中..." : "ログイン"}
                            </button>
        
                            <button
                                type="button"
                                className="text-sm text-gray-600 hover:underline"
                                onClick={() => router.push("/")}
                            >
                                キャンセル
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
}