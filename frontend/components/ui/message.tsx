type Message = {
    userid: string
    content: string
    isOwn?: boolean
}

export default function Message({ content, isOwn = false }: Message){
    const container = isOwn ? "flex items-start gap-3 justify-end" : "flex items-start gap-3";
    const bubble = isOwn
      ? "bg-blue-500 text-white px-3 py-2 rounded-lg shadow-sm max-w-[70%] text-sm"
      : "bg-white text-gray-800 px-3 py-2 rounded-lg shadow-sm max-w-[70%] text-sm";
    const avatar = isOwn ? "w-8 h-8 rounded-full bg-blue-400 flex-shrink-0" : "w-8 h-8 rounded-full bg-gray-300 flex-shrink-0";

    return (
      <div className={container}>
        {!isOwn && <div className={avatar} />}
        <div className={bubble}>{content}</div>
        {isOwn && <div className={avatar} />}
      </div>
    )
}
