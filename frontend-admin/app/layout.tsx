import type { ReactNode } from "react";
import "./globals.css";
import { SidebarProvider } from "../components/ui/SidebarContext";
import { LayoutContent } from "./layout-content";

export default function RootLayout({ children }: { children: ReactNode }) {
	return (
		<html lang="en">
			<body className="min-h-screen bg-zinc-50 text-zinc-900">
				<SidebarProvider>
					<LayoutContent>{children}</LayoutContent>
				</SidebarProvider>
			</body>
		</html>
	);
}
