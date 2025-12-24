"use client"

import { useEffect, useState } from "react";
import Room from "../ui/room";

const API_BASE = process.env.NEXT_PUBLIC_API_BASE ?? "http://localhost:3001";

type RoomSummary = {
  id: string;
  name: string;
};

interface FriendListProps {
  selectedRoomId: string;
  onRoomSelect: (roomId: string) => void;
}

export default function FriendList({ selectedRoomId, onRoomSelect }: FriendListProps) {
  const [roomList, setRoomList] = useState<RoomSummary[]>([]);
  const [newFriend, setNewFriend] = useState<string>("");
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchRooms();
  }, []);

  const fetchRooms = async () => {
    const userId = typeof window !== "undefined" ? localStorage.getItem("userName") : null;
    if (!userId) {
      setError("ログインしてください。");
      return;
    }
    try { 
      const res = await fetch(`${API_BASE}/rooms/user`, {
        method: "GET",
        headers: {
          "x-user-id": userId,
        },
      });
      // 204 No Content の場合は空配列を返す
      if (res.status === 204) {
        setRoomList([]);
        setError(null);
        return;
      }

      const data = await res.json().catch(() => null);

      if (!res.ok) {
        const msg = (data && (data.message || data.error)) || `エラー: ${res.status}`;
        setError(msg);
        return;
      }

      const roomsData = Array.isArray(data) ? data : data?.rooms || [];
      const rooms: RoomSummary[] = roomsData
        .map((r: any) => ({
          id: r.room_id || r.roomId || String(r.id || ""),
          name: r.room_name || r.roomName || "Room",
        }))
        .filter((r: RoomSummary) => !!r.id);

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
    const userId = typeof window !== "undefined" ? localStorage.getItem("userName") : null;
    if (!userId) {
      setError("ログインしてください。");
      return;
    }
    if (!newFriend.trim()) {
      setError("相手のユーザー名を入力してください。");
      return;
    }
    setNewFriend("");
    setLoading(true);
    try{
      const res = await fetch(`${API_BASE}/user/rooms`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "x-user-id": userId,
        },
        body: JSON.stringify({ user_name: newFriend.trim() }),
      });
      const data = await res.json().catch(() => null);
      if(!res.ok){
        const msg = (data && (data.message || data.error)) || `エラー: ${res.status}`;
        setError(msg);
        setLoading(false);
        return;
      }
      if (data) {
        const newRoom: RoomSummary = {
          id: data.room_id || data.roomId || String(data.id || newFriend),
          name: data.room_name || data.roomName || newFriend,
        };
        setRoomList((prev) => {
          // 既存の同じIDがあれば追加しない
          if (prev.some((r) => r.id === newRoom.id)) return prev;
          return [...prev, newRoom];
        });
      }

      setSuccess("フレンドを追加しました。");

    }catch(err){
      setError("ネットワークエラーが発生しました。");
    } finally{
      setLoading(false);
    }
  }
  
  return (
    <div className="w-80 h-full bg-gray-100 flex flex-col shadow-md">
      <ul className="flex-auto overflow-y-auto">
        {roomList.map((room) => (
          <Room 
            key={room.id}
            name={room.name}
            subtitle={room.id}
            isSelected={room.id === selectedRoomId}
            onClick={() => onRoomSelect(room.id)}
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