type Props = {
    name: string;
    onClick?: () => void;
    isSelected?: boolean;
}

export default function Room({name, onClick, isSelected}: Props) {
    const baseClass = "w-full h-20 border-b flex items-center px-4 hover:bg-gray-200 cursor-pointer font-semibold transition-colors";
    const selectedClass = isSelected 
      ? "bg-blue-200 border-l-4 border-l-blue-500" 
      : "bg-white";
    
    return (
        <li 
          className={`${baseClass} ${selectedClass}`}
          onClick={onClick}
        >
          <span className="rounded-full bg-red-300 w-10 h-10"/>
          <div className="p-4">{name}</div>
        </li>
    )
}
