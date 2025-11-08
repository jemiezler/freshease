"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { motion, AnimatePresence } from "framer-motion";
import {
  LayoutDashboard,
  Box,
  Users,
  ShoppingCart,
  MapPin,
  Bell,
  Star,
  Gift,
  Calendar,
  ChefHat,
  Truck,
  BarChart3,
  Receipt,
  CreditCard,
  ChevronDown,
} from "lucide-react";
import type { LucideIcon } from "lucide-react";
import { useState } from "react";

type NavLink = {
  href: string;
  label: string;
  icon: LucideIcon;
  children?: Array<{ href: string; label: string }>;
};

const links: NavLink[] = [
  { href: "/", label: "Dashboard", icon: LayoutDashboard },
  {
    href: "/crm",
    label: "CRM",
    icon: BarChart3,
    children: [
      { href: "/crm", label: "CRM Dashboard" },
      { href: "/crm/customers", label: "Customers" },
      { href: "/crm/orders", label: "Orders" },
      { href: "/crm/analytics", label: "Analytics" },
    ],
  },
  {
    href: "/products",
    label: "Products",
    icon: Box,
    children: [
      { href: "/categories", label: "Categories" },
      { href: "/products", label: "Products" },
      { href: "/inventories", label: "Inventories" },
      { href: "/vendors", label: "Vendors" },
    ],
  },
  { href: "/bundles", label: "Bundles", icon: Gift },
  {
    href: "/users",
    label: "Users",
    icon: Users,
    children: [
      { href: "/users", label: "Users" },
      { href: "/roles", label: "Roles" },
      { href: "/permissions", label: "Permissions" },
    ],
  },
  { href: "/carts", label: "Carts", icon: ShoppingCart, children: [
    { href: "/carts", label: "Carts" },
    { href: "/cart-items", label: "Cart Items" },
  ] },
  { href: "/orders", label: "Orders", icon: Receipt, children: [
    { href: "/orders", label: "Orders" },
    { href: "/order-items", label: "Order Items" },
  ] },
  { href: "/payments", label: "Payments", icon: CreditCard, children: [
    { href: "/payments", label: "Payments" },
    { href: "/payment-items", label: "Payment Items" },
  ] },
  { href: "/addresses", label: "Addresses", icon: MapPin, children: [
    { href: "/addresses", label: "Addresses" },
    { href: "/address-items", label: "Address Items" },
  ] },
  { href: "/deliveries", label: "Deliveries", icon: Truck, children: [
    { href: "/deliveries", label: "Deliveries" },
    { href: "/delivery-items", label: "Delivery Items" },
  ] },
  { href: "/notifications", label: "Notifications", icon: Bell },
  { href: "/reviews", label: "Reviews", icon: Star, children: [
    { href: "/reviews", label: "Reviews" },
  ] },
  { href: "/meal-plans", label: "Meal Plans", icon: Calendar, children: [
    { href: "/meal-plans", label: "Meal Plans" },
    { href: "/meal-plan-items", label: "Meal Plan Items" },
  ] },
  { href: "/recipes", label: "Recipes", icon: ChefHat, children: [
    { href: "/recipes", label: "Recipes" },
    { href: "/recipe-items", label: "Recipe Items" },
  ] },
];

const itemVariants = {
  closed: { opacity: 0, x: -20 },
  open: { opacity: 1, x: 0 },
};

const childItemVariants = {
  closed: { opacity: 0, x: -10, height: 0 },
  open: { opacity: 1, x: 0, height: "auto" },
};

const sidebarVariants = {
  open: {
    transition: { staggerChildren: 0, delayChildren: 0.1 },
  },
};

export function Sidebar() {
  const pathname = usePathname();
  const [expandedItems, setExpandedItems] = useState<Set<string>>(
    new Set(
      links
        .filter((l) => {
          if (!l.children) return false;
          return l.children.some(
            (child) =>
              pathname === child.href || pathname?.startsWith(child.href + "/")
          );
        })
        .map((l) => l.href)
    )
  );

  const toggleExpanded = (href: string) => {
    setExpandedItems((prev) => {
      const next = new Set(prev);
      if (next.has(href)) {
        next.delete(href);
      } else {
        next.add(href);
      }
      return next;
    });
  };

  return (
    <motion.aside
      initial={{ x: -220 }}
      animate={{ x: 0 }}
      transition={{ type: "spring", damping: 25, stiffness: 200 }}
      className="fixed inset-y-0 left-0 z-40 w-[220px] border-r bg-white/95 backdrop-blur supports-backdrop-filter:bg-white/70 overflow-y-auto"
    >
      <div className="px-3 py-3">
        <motion.div
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          className="mb-3 flex items-center gap-2 px-2 text-sm font-semibold text-zinc-800"
        >
          <motion.div
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ delay: 0.2, type: "spring" }}
            className="h-6 w-6 rounded bg-zinc-900"
          />
          <span>Freshease Admin</span>
        </motion.div>
        <motion.nav
          variants={sidebarVariants}
          initial="closed"
          animate="open"
          className="grid gap-1"
        >
          {links.map((l) => {
            const isActive =
              pathname === l.href ||
              (l.children && pathname?.startsWith(l.href));
            const Icon = l.icon;
            const isExpanded = l.children ? expandedItems.has(l.href) : false;
            const hasActiveChild =
              l.children &&
              l.children.some(
                (child) =>
                  pathname === child.href ||
                  pathname?.startsWith(child.href + "/")
              );

            if (l.children) {
              return (
                <motion.div
                  key={l.href}
                  variants={itemVariants}
                  className="space-y-1"
                >
                  <motion.div
                    whileHover={{ scale: 1.0, x: 2 }}
                    whileTap={{ scale: 1 }}
                  >
                    <Link
                      href={l.href}
                      onClick={(e) => {
                        e.preventDefault();
                        toggleExpanded(l.href);
                      }}
                      className={[
                        "flex items-center justify-between gap-2 rounded-md px-2 py-2 text-sm cursor-pointer transition-colors",
                        isActive || hasActiveChild
                          ? "bg-zinc-900 text-white"
                          : "text-zinc-700 hover:bg-zinc-100",
                      ].join(" ")}
                    >
                      <div className="flex items-center gap-2">
                        <Icon className="h-4 w-4" />
                        <span className="truncate">{l.label}</span>
                      </div>
                      <motion.div
                        animate={{ rotate: isExpanded || hasActiveChild ? 180 : 0 }}
                        transition={{ duration: 0.2 }}
                      >
                        <ChevronDown className="h-3 w-3" />
                      </motion.div>
                    </Link>
                  </motion.div>
                  <AnimatePresence>
                    {(isExpanded || hasActiveChild) && (
                      <motion.div
                        initial="closed"
                        animate="open"
                        exit="closed"
                        variants={{
                          open: {
                            transition: { staggerChildren: 0.03 },
                          },
                          closed: {
                            transition: { staggerChildren: 0.02, staggerDirection: -1 },
                          },
                        }}
                        className="ml-4 space-y-1 border-l-2 border-zinc-200 pl-2 overflow-hidden"
                      >
                        {l.children.map((child) => {
                          const childActive =
                            pathname === child.href ||
                            pathname?.startsWith(child.href + "/");
                          return (
                            <motion.div
                              key={child.href}
                              variants={childItemVariants}
                              transition={{ duration: 0.2 }}
                            >
                              <motion.div
                                whileHover={{ x: 4 }}
                                whileTap={{ scale: 0.98 }}
                              >
                                <Link
                                  href={child.href}
                                  className={[
                                    "flex items-center gap-2 rounded-md px-2 py-1.5 text-xs transition-colors",
                                    childActive
                                      ? "bg-zinc-100 text-zinc-900 font-medium"
                                      : "text-zinc-600 hover:bg-zinc-50",
                                  ].join(" ")}
                                >
                                  <motion.span
                                    initial={{ width: 0 }}
                                    animate={{
                                      width: childActive ? 2 : 0,
                                    }}
                                    className="h-full bg-zinc-900 rounded-full"
                                  />
                                  <span className="truncate">{child.label}</span>
                                </Link>
                              </motion.div>
                            </motion.div>
                          );
                        })}
                      </motion.div>
                    )}
                  </AnimatePresence>
                </motion.div>
              );
            }

            return (
              <motion.div key={l.href} variants={itemVariants}>
                <motion.div
                  whileHover={{ scale: 1.0, x: 4 }}
                  whileTap={{ scale: 1 }}
                >
                  <Link
                    href={l.href}
                    className={[
                      "flex items-center gap-2 rounded-md px-2 py-2 text-sm transition-colors relative",
                      isActive
                        ? "bg-zinc-900 text-white"
                        : "text-zinc-700 hover:bg-zinc-100",
                    ].join(" ")}
                  >

                    <Icon className="h-4 w-4" />
                    <span className="truncate">{l.label}</span>
                  </Link>
                </motion.div>
              </motion.div>
            );
          })}
        </motion.nav>
      </div>
    </motion.aside>
  );
}
