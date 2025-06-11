// ThemeSwitcher.tsx
import { useEffect } from 'react';
import { useTheme } from '../../hooks/useTheme';
import { themes } from '~/providers/ThemeProvider';
import ThemeOption from './ThemeOption';

export function ThemeSwitcher() {
    const { theme, setTheme } = useTheme();


    useEffect(() => {
        const saved = localStorage.getItem('theme');
        if (saved && saved !== theme) {
            setTheme(saved);
        } else if (!saved) {
            localStorage.setItem('theme', theme)
        }
    }, []);

    useEffect(() => {
        if (theme) {
            localStorage.setItem('theme', theme)
        }
    }, [theme]);

    return (
        <>
        <h2>Select Theme</h2>
        <div className="grid grid-cols-2 items-center gap-2">
            {themes.map((t) => (
                <ThemeOption setTheme={setTheme} key={t.name} theme={t} />
            ))}
        </div>
        </>
    );
}
