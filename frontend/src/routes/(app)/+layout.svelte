<script lang="ts">
    import '../../app.css';
    import favicon from '$lib/assets/favicon.svg';
    import {goto} from '$app/navigation';
    import {CircleCheckBig, Clipboard, HelpCircle, House, Plus, Settings, SquarePen} from 'lucide-svelte';

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

<div class="flex min-h-screen flex-col bg-base-100">
    <header class="sticky top-0 z-30 w-full">
        <div class="navbar border-b border-base-200 bg-base-100 shadow-lg">
            <div class="navbar-start">
                <div class="lg:hidden">
                    <button class="btn btn-ghost" onclick={toggleSidebar} aria-label="Open sidebar">
                        <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-5 w-5"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                        >
                            <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M4 6h16M4 12h8m-8 6h16"
                            />
                        </svg>
                    </button>
                </div>
                <a
                        class="text-xl font-bold text-primary normal-case"
                        onclick={navigateToHome}
                >
                    Tournament Manager
                </a>
            </div>
            <div class="navbar-end"></div>
        </div>
    </header>

    <div class="relative flex flex-grow">
        <!-- Sidebar -->
        <aside
                class={`fixed inset-y-0 z-20 w-64 shrink-0 border-r border-base-300 bg-base-200 pt-8 transition-all duration-300 ease-in-out ${sidebarOpen ? 'left-0' : '-left-64'} h-screen overflow-y-auto lg:static lg:left-0 lg:h-auto`}
        >
            <div class="p-4">
                <!-- Sidebar Header with close button on mobile -->
                <div class="mb-6 flex items-center justify-between lg:hidden">
                    <h2 class="text-lg font-semibold">Menu</h2>
                    <button class="btn btn-ghost btn-sm" onclick={toggleSidebar} aria-label="Open sidebar">
                        <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-5 w-5"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                        >
                            <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M6 18L18 6M6 6l12 12"
                            />
                        </svg>
                    </button>
                </div>

                <div class="flex flex-col space-y-1">
                    <div
                            class="mb-2 px-2 text-xs font-semibold tracking-wider text-base-content/60 uppercase"
                    >
                        Main
                    </div>
                    <a
                            href="/"
                            class="flex items-center rounded-lg px-2 py-2 text-base-content hover:bg-base-300 active:bg-primary active:text-primary-content"
                    >
                        <House class="mr-3 h-5 w-5"/>
                        Dashboard
                    </a>
                    <a
                            href="/tournaments/create"
                            class="flex items-center rounded-lg px-2 py-2 text-base-content hover:bg-base-300"
                    >
                        <Plus class="mr-3 h-5 w-5"/>
                        Create Tournament
                    </a>
                </div>

                <div class="mt-6 flex flex-col space-y-1">
                    <div
                            class="mb-2 px-2 text-xs font-semibold tracking-wider text-base-content/60 uppercase"
                    >
                        Tournaments
                    </div>
                    <a
                            href="/"
                            class="flex items-center rounded-lg px-2 py-2 text-base-content hover:bg-base-300"
                    >
                        <Clipboard class="mr-3 h-5 w-5"/>
                        All Tournaments
                    </a>
                    <a
                            href="/"
                            class="flex items-center rounded-lg px-2 py-2 text-base-content hover:bg-base-300"
                    >
                        <CircleCheckBig class="mr-3 h-5 w-5"/>
                        Active Tournaments
                    </a>
                    <a
                            href="/"
                            class="flex items-center rounded-lg px-2 py-2 text-base-content hover:bg-base-300"
                    >
                        <SquarePen class="mr-3 h-5 w-5"/>
                        Draft Tournaments
                    </a>
                </div>

                <div class="mt-6 flex flex-col space-y-1">
                    <div
                            class="mb-2 px-2 text-xs font-semibold tracking-wider text-base-content/60 uppercase"
                    >
                        Settings
                    </div>
                    <a
                            href="#"
                            class="flex items-center rounded-lg px-2 py-2 text-base-content hover:bg-base-300"
                    >
                        <Settings class="mr-3 h-5 w-5"/>
                        Settings
                    </a>
                    <a
                            href="#"
                            class="flex items-center rounded-lg px-2 py-2 text-base-content hover:bg-base-300"
                    >
                        <HelpCircle class="mr-3 h-5 w-5"/>
                        Help
                    </a>
                </div>
            </div>
        </aside>

        <!-- Overlay for mobile sidebar -->
        {#if sidebarOpen}
            <button
                    aria-label="Close sidebar"
                    class="bg-opacity-50 fixed inset-0 z-10 bg-black lg:hidden"
                    onclick={toggleSidebar}
            ></button>
        {/if}

        <!-- Main content -->
        <main class="w-full flex-grow overflow-x-hidden p-4 md:p-8 lg:pl-8">
            {@render children?.()}
        </main>
    </div>

    <footer class="footer-center footer border-t border-base-300 bg-base-200 p-6 text-base-content">
        <div class="grid grid-flow-col gap-4">
            <a href="/" class="link link-hover">Impressum</a>
            <a href="/" class="link link-hover">Datenschutz</a>
        </div>
        <div>
            <p>© {currentYear} Tournament Manager - All rights reserved</p>
        </div>
    </footer>
</div>
