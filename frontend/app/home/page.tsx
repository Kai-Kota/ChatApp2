import ChatField from "@/components/layout/chats";
import Friends from "@/components/layout/friends";
import Header from "@/components/layout/header";

export default function Home() {
  return (
    <div className="w-full h-screen flex flex-col">
      <Header/>
      <div className="flex flex-1">
        <Friends/>
        <ChatField/>
      </div>
    </div>
  );
}