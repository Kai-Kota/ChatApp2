"use client"

import { useEffect, useState, useRef, use } from "react";
import Message from "../ui/message";
import { MessageType } from "../../type"; 

export default function ChatField() {
  const [messageList, setMessageList] = useState<MessageType[]>([]);
  const [newMessage, setNewMessage] = useState<string>("");
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    fetchMessages();

    const websocket = new WebSocket("ws://localhost:8080/ws/messages");
    socketRef.current = websocket;

    const onMessage = (e: MessageEvent) => {
      const data = JSON.parse(e.data);
      let data2ws: MessageType[] = [];
      data2ws.push({
        userId: data.UserID,
        content: data.Content,
      });
      setMessageList(messageList => [...messageList, ...data2ws]);
    }
    websocket.addEventListener("message", onMessage);

    return () => {
      websocket.removeEventListener("message", onMessage);
      websocket.close();
    };
  }, [])

  const fetchMessages = async () => {
    try {
      const res = await fetch("http://localhost:8080/user/rooms", {
        method: "GET",
        credentials: "include",
         headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          room_id: 21, 
        })
      });

      let data = await res.json();
      if (!res.ok) {
        console.log("fetchError");
      }

      let messages: MessageType[] = [];
        for(let i = 0; i < data.data.length; i++){
          messages.push({
            userId: data.data[i].UserID,
            content: data.data[i].Content,
          });
        }
        setMessageList(messages);      
    } catch (err) {
      console.log("networkError");
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    //websocketで送信
    e.preventDefault();
    const data = JSON.stringify({
      content: newMessage,
      room_id: 21,
    });

    socketRef.current?.send(data);
    setNewMessage("");
  }
  return (
    <div className="flex flex-col w-full h-full">
      <div className="bg-blue-100 flex-1 h-130 p-4 space-y-3">
      {messageList.map((msg, index) => (
        <Message
          key={msg + "_" + index}
          userid={msg.userId}
          content={msg.content}  
        />
        ))}
      </div>
      <div className="p-4 bg-white flex items-end gap-2">
        <form onSubmit={handleSubmit} className="flex w-full gap-2">
          <textarea
            className="flex-1 min-h-[100px] max-h-40 resize-none rounded-md px-3 py-2 text-sm shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-300"
            onChange={(e) => setNewMessage(e.target.value)}
            placeholder="メッセージを入力..."
            aria-label="メッセージ入力"
          />
          <button 
            type = "submit"
            className="bg-blue-500 text-white px-4 py-2 rounded-md hover:bg-blue-600 focus:outline-none"
          >
            送信
          </button>
        </form>
      </div>
    </div>
  );
}