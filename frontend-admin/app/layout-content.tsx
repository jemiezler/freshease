"use client";

import type { ReactNode } from "react";
import { Sidebar } from "../components/ui/Sidebar";
import { Topbar } from "../components/ui/Topbar";
import { useSidebar } from "../components/ui/SidebarContext";
import { motion } from "framer-motion";

export function LayoutContent({ children }: { children: ReactNode }) {
	const { isCollapsed } = useSidebar();
	
	return (
		<div className="flex">
			<Sidebar />
			<motion.div 
				className="min-h-screen w-full"
				animate={{ marginLeft: isCollapsed ? 64 : 220 }}
				transition={{ type: "spring", damping: 25, stiffness: 200 }}
			>
				<Topbar />
				<main className="p-4 sm:p-6">
					<div className="mx-auto w-full max-w-[1600px]">
						{children}
					</div>
				</main>
			</motion.div>
		</div>
	);
}

