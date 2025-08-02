<script lang="ts">
    import '../app.css';
    import favicon from '$lib/assets/favicon.svg';
    import {goto} from '$app/navigation';
    import {CheckCircle, CircleCheckBig, Clipboard, HelpCircle, House, Plus, Settings, SquarePen} from 'lucide-svelte';

    let {children} = $props();

    // Sidebar state
    let sidebarOpen = true;

    // Toggle sidebar on mobile
    function toggleSidebar() {
        sidebarOpen = !sidebarOpen;
    }

    function navigateToHome() {
        goto('/');
    }

    function navigateToCreate() {
        goto('/tournaments/create');
    }

    // Get current year for footer
    const currentYear = new Date().getFullYear();
</script>

<svelte:head>
    <link rel="icon" href={favicon}/>
    <title>Tournament Manager</title>
    <meta name="description" content="Professional tournament management system"/>
</svelte:head>

<div class="min-h-screen flex flex-col bg-base-100">
    <header class="sticky top-0 z-30 w-full">
        <div class="navbar bg-base-100 shadow-lg border-b border-base-200">
            <div class="navbar-start">
                <div class="lg:hidden">
                    <button class="btn btn-ghost" onclick={toggleSidebar} aria-label="Open sidebar">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 24 24"
                             stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M4 6h16M4 12h8m-8 6h16"/>
                        </svg>
                    </button>
                </div>
                <button class="btn btn-ghost normal-case text-xl font-bold text-primary" onclick={navigateToHome}>
                    Tournament Manager
                </button>
            </div>
            <div class="navbar-end">

            </div>
        </div>
    </header>

    <div class="flex flex-grow relative">
        <!-- Sidebar -->
        <aside class={`bg-base-200 w-64 shrink-0 border-r border-base-300 fixed inset-y-0 pt-8 z-20 transition-all duration-300 ease-in-out ${sidebarOpen ? 'left-0' : '-left-64'} lg:left-0 lg:static h-screen lg:h-auto overflow-y-auto`}>
            <div class="p-4">
                <!-- Sidebar Header with close button on mobile -->
                <div class="flex items-center justify-between mb-6 lg:hidden">
                    <h2 class="text-lg font-semibold">Menu</h2>
                    <button class="btn btn-sm btn-ghost" onclick={toggleSidebar} aria-label="Open sidebar">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 24 24"
                             stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                    </button>
                </div>

                <div class="flex flex-col space-y-1">
                    <div class="text-xs font-semibold text-base-content/60 uppercase tracking-wider mb-2 px-2">
                        Main
                    </div>
                    <a href="/"
                       class="flex items-center px-2 py-2 text-base-content rounded-lg hover:bg-base-300 active:bg-primary active:text-primary-content">
                        <House class="h-5 w-5 mr-3"/>
                        Dashboard
                    </a>
                    <a href="/tournaments/create"
                       class="flex items-center px-2 py-2 text-base-content rounded-lg hover:bg-base-300">
                        <Plus class="h-5 w-5 mr-3"/>
                        Create Tournament
                    </a>
                </div>

                <div class="flex flex-col space-y-1 mt-6">
                    <div class="text-xs font-semibold text-base-content/60 uppercase tracking-wider mb-2 px-2">
                        Tournaments
                    </div>
                    <a href="/" class="flex items-center px-2 py-2 text-base-content rounded-lg hover:bg-base-300">
                        <Clipboard class="h-5 w-5 mr-3"/>
                        All Tournaments
                    </a>
                    <a href="/" class="flex items-center px-2 py-2 text-base-content rounded-lg hover:bg-base-300">
                        <CircleCheckBig class="h-5 w-5 mr-3"/>
                        Active Tournaments
                    </a>
                    <a href="/" class="flex items-center px-2 py-2 text-base-content rounded-lg hover:bg-base-300">
                        <SquarePen class="h-5 w-5 mr-3"/>
                        Draft Tournaments
                    </a>
                </div>

                <div class="flex flex-col space-y-1 mt-6">
                    <div class="text-xs font-semibold text-base-content/60 uppercase tracking-wider mb-2 px-2">
                        Settings
                    </div>
                    <a href="#" class="flex items-center px-2 py-2 text-base-content rounded-lg hover:bg-base-300">
                        <Settings class="h-5 w-5 mr-3"/>
                        Settings
                    </a>
                    <a href="#" class="flex items-center px-2 py-2 text-base-content rounded-lg hover:bg-base-300">
                        <HelpCircle class="h-5 w-5 mr-3"/>
                        Help
                    </a>
                </div>
            </div>
        </aside>

        <!-- Overlay for mobile sidebar -->
        {#if sidebarOpen}
            <button
                    aria-label="Close sidebar"
                    class="fixed inset-0 bg-black bg-opacity-50 z-10 lg:hidden"
                    onclick={toggleSidebar}
            ></button>
        {/if}

        <!-- Main content -->
        <main class="flex-grow p-4 md:p-8 lg:pl-8 w-full overflow-x-hidden">
            {@render children?.()}
        </main>
    </div>

    <footer class="footer footer-center p-6 bg-base-200 text-base-content border-t border-base-300">
        <div class="grid grid-flow-col gap-4">
            <a href="/" class="link link-hover">Impressum</a>
            <a href="/" class="link link-hover">Datenschutz</a>
        </div>
        <div>
            <p>© {currentYear} Tournament Manager - All rights reserved</p>
        </div>
    </footer>
</div>
