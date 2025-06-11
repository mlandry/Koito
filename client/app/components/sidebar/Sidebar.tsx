import { ExternalLink, Home,  Info } from "lucide-react";
import SidebarSearch from "./SidebarSearch";
import SidebarItem from "./SidebarItem";
import SidebarSettings from "./SidebarSettings";

export default function Sidebar() {

    const iconSize = 20;

    return (
        <div className="z-50 flex flex-col justify-between h-screen border-r-1 border-(--color-bg-tertiary) p-1 py-10 sticky left-0 top-0 bg-(--color-bg)">
            <div className="flex flex-col gap-4">
                <SidebarItem space={10} to="/" name="Home" onClick={() => {}} modal={<></>}><Home size={iconSize} /></SidebarItem>
                <SidebarSearch size={iconSize} />
            </div>
            <div className="flex flex-col gap-4">
                <SidebarItem icon keyHint={<ExternalLink size={14} />} space={22} externalLink to="https://koito.io" name="About" onClick={() => {}} modal={<></>}><Info size={iconSize} /></SidebarItem>
                <SidebarSettings size={iconSize} />
            </div>
        </div>
    );
}