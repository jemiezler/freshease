"use client";

import type { ReactNode } from "react";
import { usePathname } from "next/navigation";
import { Sidebar } from "../components/ui/Sidebar";
import { Topbar } from "../components/ui/Topbar";
import { useSidebar } from "../components/ui/SidebarContext";
import { useAuth } from "../lib/auth-context";
import { motion } from "framer-motion";

export function LayoutContent({ children }: { children: ReactNode }) {
	const { isCollapsed } = useSidebar();
	const pathname = usePathname();
	const { loading, isAuthenticated } = useAuth();
	
	const isLoginPage = pathname === "/login";
	const isInitAdminPage = pathname === "/init-admin";
	
	// Show login/init-admin page without sidebar/topbar
	if (isLoginPage || isInitAdminPage) {
		return <>{children}</>;
	}
	
	// Show loading state while checking auth
	if (loading) {
		return (
			<div className="flex min-h-screen items-center justify-center">
				<div className="text-center">
					<div className="mb-4 inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-zinc-900 border-r-transparent"></div>
					<p className="text-sm text-zinc-600">Loading...</p>
				</div>
			</div>
		);
	}
	
	// Show admin layout for authenticated users
	if (isAuthenticated) {
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
	
	// Not authenticated, show nothing (redirect will happen in AuthProvider)
	return null;
}

