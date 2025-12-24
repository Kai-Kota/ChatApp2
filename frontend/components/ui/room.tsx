type Props = {
  name: string;
  subtitle?: string;
  isSelected?: boolean;
  onClick?: () => void;
}

export default function Room({ name, subtitle, isSelected = false, onClick }: Props) {
    const base = "w-full h-20 border-b flex items-center px-4 cursor-pointer font-semibold";
    const selected = "bg-blue-100 border-l-4 border-blue-500";
    const unselected = "bg-white hover:bg-gray-200";
    const className = `${base} ${isSelected ? selected : unselected}`;

    return (
        <li
          className={className}
          onClick={onClick}
        >
          <span className="rounded-full bg-red-300 w-10 h-10"/>
          <div className="p-4">
            <div>{name}</div>
            {subtitle && (
              <div className="text-xs text-gray-500">{subtitle}</div>
            )}
          </div>
        </li>
    )
}