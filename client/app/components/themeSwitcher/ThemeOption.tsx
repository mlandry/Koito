import type { Theme } from "~/providers/ThemeProvider";

interface Props {
    theme: Theme
    setTheme: Function
}

export default function ThemeOption({ theme, setTheme }: Props) {

    const capitalizeFirstLetter = (s: string) => {
        return s.charAt(0).toUpperCase() + s.slice(1);
    }

    return (
        <div onClick={() => setTheme(theme.name)} className="rounded-md p-5 hover:cursor-pointer flex gap-4 items-center border-2" style={{background: theme.bg, color: theme.fg, borderColor: theme.bgSecondary}}>
            {capitalizeFirstLetter(theme.name)}
            <div className="w-[50px] h-[30px] rounded-md" style={{background: theme.bgSecondary}}></div>
            <div className="w-[50px] h-[30px] rounded-md" style={{background: theme.fgSecondary}}></div>
            <div className="w-[50px] h-[30px] rounded-md" style={{background: theme.primary}}></div>
        </div>
    )
}