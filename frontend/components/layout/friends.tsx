"use client"

import { useEffect, useState } from "react";
import Room from "../ui/room";

type Room = {
  id: number
}

export default function FriendList() {
  const [roomList, setRoomList] = useState<Room[]>([]);
  const [newFriend, setNewFriend] = useState<string>("");
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchRooms();
  }, []);

  const fetchRooms = async () => {
    try { 
      const res = await fetch("http://localhost:8080/user/rooms", {
        method: "GET",
        credentials: "include",
      });
      // 204 No Content の場合は空配列を返す
      if (res.status === 204) {
        setRoomList([]);
        setError(null);
        return;
      }

      // Content-Type を見て JSON かどうかを判定し、安全に読み取る
      const contentType = (res.headers.get("content-type") || "").toLowerCase();
      let data: any = null;
      if (contentType.includes("application/json")) {
        try {
          data = await res.json();
        } catch (e) {
          data = null;
        }
      } else {
        // JSON でない場合は text として読み取り、JSON なら parse する
        const text = await res.text();
        try {
          data = text ? JSON.parse(text) : null;
        } catch (e) {
          data = text || null;
        }
      }

      if (!res.ok) {
        const msg = (data && (data.message || data.error)) || `エラー: ${res.status}`;
        setError(msg);
      }

      let rooms: Room[] = [];
      for(let i = 0; i < data.data.length; i++){;
        rooms.push({id: data.data[i].ID})
      }
      setRoomList(rooms);
      setError(null);
    } catch (err) {
      setError("ネットワークエラーが発生しました。");
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);
    setNewFriend("");
    setLoading(true);
    try{
      const res = await fetch("http://localhost:8080/user/rooms", {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ pertner: newFriend }),
      });
      if(!res.ok){
        const data = await res.json().catch(() => ({}));
        setError(data.message || `エラー: ${res.status}`);
        setLoading(false);
        return;
      }
      if(res.ok) console.log("Friend added");

      setSuccess("フレンドを追加しました。");

    }catch(err){
      setError("ネットワークエラーが発生しました。");
    } finally{
      setLoading(false);
    }
  }
  
  return (
    <div className="w-80 h-full bg-gray-100 flex flex-col shadow-md">
      <ul className="flex-auto">
        {roomList.map((room) => (
          <Room 
            key={room.id}
            name={`Friend ${room.id}`}
          />
        ))}
      </ul>

      <form onSubmit={handleSubmit} className="w-full h-25 bg-cyan-800 p-3">
        <p className="text-white">フレンド追加</p>
        <div className="flex justify-between">
          <input 
            type="text"
            value={newFriend}
            onChange={(e) => setNewFriend(e.target.value)}
            className="mt-2 p-2 rounded-md border bg-white border-gray-300 text-sm"
          />
          <input 
            type="submit"
            value="追加"
            className=" mt-2 p-2 rounded-md bg-blue-500 text-white text-sm hover:bg-blue-600 cursor-pointer"
           />
         </div>
      </form>
    </div>
  );
}