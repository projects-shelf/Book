import { SidebarProvider, SidebarTrigger, useSidebar } from "@/components/ui/sidebar";
import { AppSidebar } from "./components/app-sidebar";
import { Outlet } from "react-router-dom";

function Trigger() {
    const { open, openMobile, isMobile } = useSidebar();

    if (isMobile ? openMobile : open) {
        return null;
    }

    return <SidebarTrigger className="fixed h-9 w-9 [&_svg]:size-7" />
}

export default function Layout() {
    return (
        <SidebarProvider>
            <AppSidebar />
            <Trigger />
            <Outlet />
        </SidebarProvider >
    );
}
