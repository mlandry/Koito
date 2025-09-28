import { globalStyle } from "@vanilla-extract/css"
import { themeVars } from "./vars.css"

export type Theme = {
    bg: string
    bgSecondary: string
    bgTertiary: string
    fg: string
    fgSecondary: string
    fgTertiary: string
    primary: string
    primaryDim: string
    accent: string
    accentDim: string
    error: string
    warning: string
    info: string
    success: string
}

export const THEME_KEYS = [
    '--color'
]

export const themes: Record<string, Theme> = {
    yuu: {
        bg: "#1e1816",
        bgSecondary: "#2f2623",
        bgTertiary: "#453733",
        fg: "#f8f3ec",
        fgSecondary: "#d6ccc2",
        fgTertiary: "#b4a89c",
        primary: "#fc9174",
        primaryDim: "#d88b65",
        accent: "#f9db6d",
        accentDim: "#d9bc55",
        error: "#e26c6a",
        warning: "#f5b851",
        success: "#8fc48f",
        info: "#87b8dd",
    },
    varia: {
        bg: "rgb(25, 25, 29)",
        bgSecondary: "#222222",
        bgTertiary: "#333333",
        fg: "#eeeeee",
        fgSecondary: "#aaaaaa",
        fgTertiary: "#888888",
        primary: "rgb(203, 110, 240)",
        primaryDim: "#c28379",
        accent: "#f0ad0a",
        accentDim: "#d08d08",
        error: "#f44336",
        warning: "#ff9800",
        success: "#4caf50",
        info: "#2196f3",
    },
    midnight: {
        bg: "rgb(8, 15, 24)",
        bgSecondary: "rgb(15, 27, 46)",
        bgTertiary: "rgb(15, 41, 70)",
        fg: "#dbdfe7",
        fgSecondary: "#9ea3a8",
        fgTertiary: "#74787c",
        primary: "#1a97eb",
        primaryDim: "#2680aa",
        accent: "#f0ad0a",
        accentDim: "#d08d08",
        error: "#f44336",
        warning: "#ff9800",
        success: "#4caf50",
        info: "#2196f3",
    },
    catppuccin: {
        bg: "#1e1e2e",
        bgSecondary: "#181825",
        bgTertiary: "#11111b",
        fg: "#cdd6f4",
        fgSecondary: "#a6adc8",
        fgTertiary: "#9399b2",
        primary: "#89b4fa",
        primaryDim: "#739df0",
        accent: "#f38ba8",
        accentDim: "#d67b94",
        error: "#f38ba8",
        warning: "#f9e2af",
        success: "#a6e3a1",
        info: "#89dceb",
    },
    autumn: {
        bg: "rgb(44, 25, 18)",
        bgSecondary: "rgb(70, 40, 18)",
        bgTertiary: "#4b2f1c",
        fg: "#fef9f3",
        fgSecondary: "#dbc6b0",
        fgTertiary: "#a3917a",
        primary: "#d97706",
        primaryDim: "#b45309",
        accent: "#8c4c28",
        accentDim: "#6b3b1f",
        error: "#d1433f",
        warning: "#e38b29",
        success: "#6b8e23",
        info: "#c084fc",
    },
    black: {
        bg: "#000000",
        bgSecondary: "#1a1a1a",
        bgTertiary: "#2a2a2a",
        fg: "#dddddd",
        fgSecondary: "#aaaaaa",
        fgTertiary: "#888888",
        primary: "#08c08c",
        primaryDim: "#08c08c",
        accent: "#f0ad0a",
        accentDim: "#d08d08",
        error: "#f44336",
        warning: "#ff9800",
        success: "#4caf50",
        info: "#2196f3",
    },
    wine: {
        bg: "#23181E",
        bgSecondary: "#2C1C25",
        bgTertiary: "#422A37",
        fg: "#FCE0B3",
        fgSecondary: "#C7AC81",
        fgTertiary: "#A78E64",
        primary: "#EA8A64",
        primaryDim: "#BD7255",
        accent: "#FAE99B",
        accentDim: "#C6B464",
        error: "#fca5a5",
        warning: "#fde68a",
        success: "#bbf7d0",
        info: "#bae6fd",
    },
    pearl: {
        bg: "#FFFFFF",
        bgSecondary: "#EEEEEE",
        bgTertiary: "#E0E0E0",
        fg: "#333333",
        fgSecondary: "#555555",
        fgTertiary: "#777777",
        primary: "#007BFF",
        primaryDim: "#0056B3",
        accent: "#28A745",
        accentDim: "#1E7E34",
        error: "#DC3545",
        warning: "#FFC107",
        success: "#28A745",
        info: "#17A2B8",
    },
    asuka: {
        bg: "#3B1212",
        bgSecondary: "#471B1B",
        bgTertiary: "#020202",
        fg: "#F1E9E6",
        fgSecondary: "#CCB6AE",
        fgTertiary: "#9F8176",
        primary: "#F1E9E6",
        primaryDim: "#CCB6AE",
        accent: "#41CE41",
        accentDim: "#3BA03B",
        error: "#DC143C",
        warning: "#FFD700",
        success: "#32CD32",
        info: "#1E90FF",
    },
    urim: {
        bg: "#101713",
        bgSecondary: "#1B2921",
        bgTertiary: "#273B30",
        fg: "#D2E79E",
        fgSecondary: "#B4DA55",
        fgTertiary: "#7E9F2A",
        primary: "#ead500",
        primaryDim: "#C1B210",
        accent: "#28A745",
        accentDim: "#1E7E34",
        error: "#EE5237",
        warning: "#FFC107",
        success: "#28A745",
        info: "#17A2B8",
    },
    match: {
        bg: "#071014",
        bgSecondary: "#0A181E",
        bgTertiary: "#112A34",
        fg: "#ebeaeb",
        fgSecondary: "#BDBDBD",
        fgTertiary: "#A2A2A2",
        primary: "#fda827",
        primaryDim: "#C78420",
        accent: "#277CFD",
        accentDim: "#1F60C1",
        error: "#F14426",
        warning: "#FFC107",
        success: "#28A745",
        info: "#17A2B8",
    },
    lemon: {
        bg: "#1a171a",
        bgSecondary: "#2E272E",
        bgTertiary: "#443844",
        fg: "#E6E2DC",
        fgSecondary: "#B2ACA1",
        fgTertiary: "#968F82",
        primary: "#f5c737",
        primaryDim: "#C29D2F",
        accent: "#277CFD",
        accentDim: "#1F60C1",
        error: "#F14426",
        warning: "#FFC107",
        success: "#28A745",
        info: "#17A2B8",
    }
};

export default themes

Object.entries(themes).forEach(([name, theme]) => {
    const selector = `[data-theme="${name}"]`

    globalStyle(selector, {
        vars: {
            [themeVars.bg]: theme.bg,
            [themeVars.bgSecondary]: theme.bgSecondary,
            [themeVars.bgTertiary]: theme.bgTertiary,
            [themeVars.fg]: theme.fg,
            [themeVars.fgSecondary]: theme.fgSecondary,
            [themeVars.fgTertiary]: theme.fgTertiary,
            [themeVars.primary]: theme.primary,
            [themeVars.primaryDim]: theme.primaryDim,
            [themeVars.accent]: theme.accent,
            [themeVars.accentDim]: theme.accentDim,
            [themeVars.error]: theme.error,
            [themeVars.warning]: theme.warning,
            [themeVars.success]: theme.success,
            [themeVars.info]: theme.info,
        }
    })
})
