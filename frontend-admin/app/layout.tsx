import type { ReactNode } from "react";
import "./globals.css";
import { Sidebar } from "../components/ui/Sidebar";
import { Topbar } from "../components/ui/Topbar";

export default function RootLayout({ children }: { children: ReactNode }) {
	return (
		<html lang="en">
			<body className="min-h-screen bg-zinc-50 text-zinc-900">
				<div className="flex">
					<Sidebar />
					<div className="ml-[220px] min-h-screen w-full">
						<Topbar />
						<main className="p-4 sm:p-6">
							<div className="mx-auto w-full max-w-[1600px]">
								{children}
							</div>
						</main>
					</div>
				</div>
			</body>
		</html>
	);
}
