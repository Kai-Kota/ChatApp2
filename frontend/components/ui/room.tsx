type Props = {
    name: string
}

export default function Room({name}: Props) {
    return (
        <li className="w-full h-20 bg-white border-b flex items-center px-4 hover:bg-gray-200 cursor-pointer font-semibold">
          <span className="rounded-full bg-red-300 w-10 h-10"/>
          <div className="p-4">{name}</div>
        </li>
    )
}