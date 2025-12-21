type Message = {
    userid: number
    content: string
}

export default function Message({content} :Message){
    return (
    // 送信用
      <div className="flex items-start gap-3 justify-end">
        <div className="bg-blue-500 text-white px-3 py-2 rounded-lg shadow-sm max-w-[70%] text-sm">
          {content}
        </div>
        <div className="w-8 h-8 rounded-full bg-blue-400 flex-shrink-0" />
      </div>     
    )
}

// 受信用
<div className="flex items-start gap-3">
  <div className="w-8 h-8 rounded-full bg-gray-300 flex-shrink-0" />
  <div className="bg-white px-3 py-2 rounded-lg shadow-sm max-w-[70%] text-sm">
    元気？
  </div>
</div>