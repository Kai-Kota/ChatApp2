"use client";

import { useState } from "react";
import ChatField from "@/components/layout/chats";
import Friends from "@/components/layout/friends";
import Header from "@/components/layout/header";

export default function Home() {
  const [selectedRoomId, setSelectedRoomId] = useState<string>("");

  return (
    <div className="w-full h-screen flex flex-col">
      <Header/>
      <div className="flex flex-1">
        <Friends selectedRoomId={selectedRoomId} onRoomSelect={setSelectedRoomId} />
        <ChatField roomId={selectedRoomId} />
      </div>
    </div>
  );
}