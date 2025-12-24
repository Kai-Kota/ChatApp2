"use client"

import { useEffect, useState, useRef } from "react";
import Message from "../ui/message";
import { MessageType } from "../../type";

const WS_URL = process.env.NEXT_PUBLIC_WS_URL ?? "ws://localhost:8080/ws";
const API_BASE = process.env.NEXT_PUBLIC_API_BASE ?? "http://localhost:3001";

interface ChatFieldProps {
  roomId?: string;
}

export default function ChatField({ roomId = "" }: ChatFieldProps) {
  const [messageList, setMessageList] = useState<MessageType[]>([]);
  const [newMessage, setNewMessage] = useState<string>("");
  const [connected, setConnected] = useState<boolean>(false);
  const [currentUserId, setCurrentUserId] = useState<string>("");
  const socketRef = useRef<WebSocket | null>(null);
  const reconnectTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const attemptsRef = useRef<number>(0);

  useEffect(() => {
    if (typeof window !== "undefined") {
      setCurrentUserId(localStorage.getItem("userName") ?? "");
    }
  }, []);

  useEffect(() => {
    // ルーム切替時にメッセージ履歴を読み込む（最新20件）
    setMessageList([]);
    const userId = typeof window !== "undefined" ? localStorage.getItem("userName") : null;
    if (roomId) {
      fetch(`${API_BASE}/rooms/${roomId}/messages?limit=20`, {
        method: "GET",
        headers: { "x-user-id": userId || "" },
      })
        .then(async (res) => {
          if (!res.ok) return [];
          const data = await res.json().catch(() => []);
          // data: [{ room_id, timestamp, user_id, content }]
          const initial = Array.isArray(data)
            ? data.map((m: any) => ({ userId: m.user_id || m.userId, content: m.content }))
            : [];
          setMessageList(initial);
        })
        .catch(() => {});
    }

    if (!userId || !roomId) {
      console.log("No userId or roomId, skipping WebSocket connection");
      return;
    }

    const connect = () => {
      const wsUrl = `${WS_URL}?room_id=${roomId}&user_id=${userId}`;
      console.log("Connecting WS:", wsUrl);
      const websocket = new WebSocket(wsUrl);
      socketRef.current = websocket;

      websocket.onopen = () => {
        attemptsRef.current = 0;
        setConnected(true);
        console.log("WebSocket connected to room:", roomId);
      };

      websocket.onmessage = (e: MessageEvent) => {
        try {
          const message = JSON.parse(e.data);
          setMessageList((prev) => [...prev, { userId: message.user_id, content: message.content }]);
        } catch (error) {
          console.error("Failed to parse message:", error);
        }
      };

      websocket.onerror = (error) => {
        console.error("WebSocket error:", error);
      };

      websocket.onclose = () => {
        setConnected(false);
        console.log("WebSocket disconnected");
        // 自動リトライ（最大5回、指数バックオフ）
        if (attemptsRef.current < 5) {
          const delay = Math.min(1000 * Math.pow(2, attemptsRef.current), 10000);
          attemptsRef.current += 1;
          reconnectTimerRef.current = setTimeout(() => {
            console.log(`Reconnecting... attempt ${attemptsRef.current}`);
            connect();
          }, delay);
        }
      };
    };

    // 初回接続
    connect();

    return () => {
      if (reconnectTimerRef.current) {
        clearTimeout(reconnectTimerRef.current);
        reconnectTimerRef.current = null;
      }
      attemptsRef.current = 0;
      if (socketRef.current) {
        socketRef.current.close();
        socketRef.current = null;
      }
    };
  }, [roomId]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!newMessage.trim() || !socketRef.current || socketRef.current.readyState !== WebSocket.OPEN) {
      return;
    }

    const messageData = JSON.stringify({
      action: "sendMessage",
      room_id: roomId,
      content: newMessage.trim(),
    });

    socketRef.current.send(messageData);
    setNewMessage("");
  };

  if (!roomId) {
    return (
      <div className="flex flex-col w-full h-full items-center justify-center bg-gray-50">
        <p className="text-gray-500">部屋を選択してください</p>
      </div>
    );
  }

  return (
    <div className="flex flex-col w-full h-full">
      <div className="bg-blue-100 h-[calc(100vh-4rem-8rem)] overflow-y-auto p-4 space-y-3">
        {!connected && roomId && (
          <div className="text-xs text-red-600">WebSocketに接続できません。再試行中...</div>
        )}
        {messageList.map((msg, index) => (
          <Message
            key={`${msg.userId}_${index}`}
            userid={msg.userId}
            content={msg.content}
            isOwn={msg.userId === currentUserId}
          />
        ))}
      </div>
      <div className="h-32 p-4 bg-white flex items-end gap-2 border-t">
        <form onSubmit={handleSubmit} className="flex w-full gap-2">
          <textarea
            value={newMessage}
            className="flex-1 h-24 resize-none rounded-md px-3 py-2 text-sm shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-300"
            onChange={(e) => setNewMessage(e.target.value)}
            placeholder="メッセージを入力..."
            aria-label="メッセージ入力"
          />
          <button type="submit" className="bg-blue-500 text-white px-4 py-2 rounded-md hover:bg-blue-600 focus:outline-none">
            送信
          </button>
        </form>
      </div>
    </div>
  );
}