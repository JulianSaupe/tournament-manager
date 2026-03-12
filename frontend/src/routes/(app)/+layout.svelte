<script lang="ts">
    import '../../app.css';
    import favicon from '$lib/assets/favicon.svg';
    import {
        CircleCheckBig,
        CircleQuestionMark,
        Clipboard,
        House,
        Menu,
        Plus,
        Settings,
        SquarePen,
        X
    } from 'lucide-svelte';

    let {children} = $props();

    // Sidebar state - closed by default on mobile
    let sidebarOpen = $state(false);

    // Get current path for active state
    let currentPath = $state('');

    $effect(() => {
        if (typeof window !== 'undefined') {
            currentPath = window.location.pathname;
        }
    });

    // Helper to check if a path is active
    function isActive(path: string): boolean {
        return currentPath === path;
    }

    // Helper to get link classes
    function getLinkClasses(path: string, isMobile: boolean = false): string {
        const baseClasses = `flex items-center gap-3 rounded-lg px-3 text-sm font-medium transition-colors hover:bg-base-200`;
        const paddingClasses = isMobile ? 'py-2.5' : 'py-2';
        const activeClasses = isActive(path) ? 'bg-primary/10 text-primary' : 'text-base-content';
        return `${baseClasses} ${paddingClasses} ${activeClasses}`;
    }

    // Close sidebar when clicking a link on mobile
    function closeSidebarOnMobile() {
        if (window.innerWidth < 1024) {
            sidebarOpen = false;
        }
    }

    // Get current year for footer
    const currentYear = new Date().getFullYear();
</script>

<svelte:head>
    <link href={favicon} rel="icon"/>
    <title>TournamentProvider Manager</title>
    <meta content="Professional tournament management system" name="description"/>
</svelte:head>

<div class="flex min-h-screen flex-col">
    <!-- Header -->
    <header class="sticky top-0 z-40 border-b border-base-300 bg-base-100 shadow-sm">
        <div class="flex h-16 items-center justify-between px-4 lg:px-6">
            <!-- Mobile menu button -->
            <button
                    aria-label="Toggle menu"
                    class="btn btn-ghost btn-sm lg:hidden"
                    onclick={() => sidebarOpen = !sidebarOpen}
            >
                <Menu class="h-5 w-5"/>
            </button>

            <!-- Logo/Title -->
            <a class="text-lg font-bold text-primary sm:text-xl" href="/">
                TournamentProvider Manager
            </a>

            <!-- Right side spacer for mobile balance -->
            <div class="w-10 lg:hidden"></div>
        </div>
    </header>

    <div class="flex flex-1">
        <!-- Sidebar for desktop -->
        <aside class="hidden w-64 border-r border-base-300 bg-base-100 lg:block">
            <nav class="flex flex-col gap-1 p-4">
                <!-- Main section -->
                <div class="mb-4">
                    <div class="mb-2 px-3 text-xs font-semibold uppercase tracking-wider text-base-content/60">
                        Main
                    </div>
                    <a class={getLinkClasses('/')} href="/">
                        <House class="h-4 w-4"/>
                        Dashboard
                    </a>
                    <a class={getLinkClasses('/tournaments/create')} href="/tournaments/create">
                        <Plus class="h-4 w-4"/>
                        Create TournamentProvider
                    </a>
                </div>

                <!-- Tournaments section -->
                <div class="mb-4">
                    <div class="mb-2 px-3 text-xs font-semibold uppercase tracking-wider text-base-content/60">
                        Tournaments
                    </div>
                    <a
                            class="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors hover:bg-base-200 text-base-content"
                            href="/"
                    >
                        <Clipboard class="h-4 w-4"/>
                        All Tournaments
                    </a>
                    <a
                            class="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors hover:bg-base-200 text-base-content"
                            href="/"
                    >
                        <CircleCheckBig class="h-4 w-4"/>
                        Active
                    </a>
                    <a
                            class="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors hover:bg-base-200 text-base-content"
                            href="/"
                    >
                        <SquarePen class="h-4 w-4"/>
                        Draft
                    </a>
                </div>

                <!-- Settings section -->
                <div class="mt-auto">
                    <div class="mb-2 px-3 text-xs font-semibold uppercase tracking-wider text-base-content/60">
                        Settings
                    </div>
                    <a
                            class="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors hover:bg-base-200 text-base-content"
                            href="#"
                    >
                        <Settings class="h-4 w-4"/>
                        Settings
                    </a>
                    <a
                            class="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors hover:bg-base-200 text-base-content"
                            href="#"
                    >
                        <CircleQuestionMark class="h-4 w-4"/>
                        Help
                    </a>
                </div>
            </nav>
        </aside>

        <!-- Mobile sidebar drawer -->
        {#if sidebarOpen}
            <!-- Backdrop -->
            <div
                    class="fixed inset-0 z-40 bg-black/50 lg:hidden"
                    onclick={() => sidebarOpen = false}
            ></div>

            <!-- Drawer -->
            <div class="fixed inset-y-0 left-0 z-50 w-64 bg-base-100 shadow-xl lg:hidden">
                <div class="flex h-16 items-center justify-between border-b border-base-300 px-4">
                    <span class="text-lg font-bold text-primary">Menu</span>
                    <button
                            class="btn btn-ghost btn-sm"
                            onclick={() => sidebarOpen = false}
                            aria-label="Close menu"
                    >
                        <X class="h-5 w-5"/>
                    </button>
                </div>

                <nav class="flex flex-col gap-1 overflow-y-auto p-4" style="height: calc(100vh - 4rem);">
                    <!-- Main section -->
                    <div class="mb-4">
                        <div class="mb-2 px-3 text-xs font-semibold uppercase tracking-wider text-base-content/60">
                            Main
                        </div>
                        <a href="/" onclick={closeSidebarOnMobile} class={getLinkClasses('/', true)}>
                            <House class="h-5 w-5"/>
                            Dashboard
                        </a>
                        <a href="/tournaments/create" onclick={closeSidebarOnMobile}
                           class={getLinkClasses('/tournaments/create', true)}>
                            <Plus class="h-5 w-5"/>
                            Create TournamentProvider
                        </a>
                    </div>

                    <!-- Tournaments section -->
                    <div class="mb-4">
                        <div class="mb-2 px-3 text-xs font-semibold uppercase tracking-wider text-base-content/60">
                            Tournaments
                        </div>
                        <a
                                href="/"
                                onclick={closeSidebarOnMobile}
                                class="flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors hover:bg-base-200 text-base-content"
                        >
                            <Clipboard class="h-5 w-5"/>
                            All Tournaments
                        </a>
                        <a
                                href="/"
                                onclick={closeSidebarOnMobile}
                                class="flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors hover:bg-base-200 text-base-content"
                        >
                            <CircleCheckBig class="h-5 w-5"/>
                            Active
                        </a>
                        <a
                                href="/"
                                onclick={closeSidebarOnMobile}
                                class="flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors hover:bg-base-200 text-base-content"
                        >
                            <SquarePen class="h-5 w-5"/>
                            Draft
                        </a>
                    </div>

                    <!-- Settings section -->
                    <div>
                        <div class="mb-2 px-3 text-xs font-semibold uppercase tracking-wider text-base-content/60">
                            Settings
                        </div>
                        <a
                                href="#"
                                onclick={closeSidebarOnMobile}
                                class="flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors hover:bg-base-200 text-base-content"
                        >
                            <Settings class="h-5 w-5"/>
                            Settings
                        </a>
                        <a
                                href="#"
                                onclick={closeSidebarOnMobile}
                                class="flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors hover:bg-base-200 text-base-content"
                        >
                            <CircleQuestionMark class="h-5 w-5"/>
                            Help
                        </a>
                    </div>
                </nav>
            </div>
        {/if}

        <!-- Main content -->
        <main class="flex-1 overflow-x-hidden p-4 md:p-6 lg:p-8">
            {@render children?.()}
        </main>
    </div>

    <!-- Footer -->
    <footer class="border-t border-base-300 bg-base-100 px-4 py-6 text-center text-sm text-base-content/70">
        <div class="mb-3 flex justify-center gap-4">
            <a class="link link-hover" href="/">Impressum</a>
            <a class="link link-hover" href="/">Datenschutz</a>
        </div>
        <p>© {currentYear} TournamentProvider Manager - All rights reserved</p>
    </footer>
</div>
