<script>
  import { Button, Card, ScrollArea, Separator } from "@nanoteck137/nano-ui";
  import NotLoggedIn from "./NotLoggedIn.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import SmallMediaCard from "$lib/components/SmallMediaCard.svelte";

  const { data } = $props();
</script>

{#if !data.user}
  <NotLoggedIn />
{:else}
  <Card.Root>
    <div class="relative w-full overflow-hidden rounded-xl shadow-lg">
      <div
        class="h-full w-full bg-gradient-to-r from-indigo-500 via-purple-500 to-pink-500 p-6"
      >
        <p class="text-sm text-muted">Welcome back</p>
        <a
          class="text-3xl font-bold hover:underline"
          href="/users/{data.user.id}"
        >
          {data.user.displayName}
        </a>
      </div>
    </div>
  </Card.Root>

  <Spacer size="md" />

  <Card.Root>
    <div class="grid grid-cols-2 justify-around gap-6 p-6 sm:grid-cols-4">
      <div class="flex flex-col items-center">
        <p class="text-xl">6</p>
        <p class="text-center text-xs">In Progress</p>
      </div>

      <div class="flex flex-col items-center">
        <p class="text-xl">6</p>
        <p class="text-center text-xs">In Progress</p>
      </div>

      <div class="flex flex-col items-center">
        <p class="text-xl">6</p>
        <p class="text-center text-xs">In Progress</p>
      </div>

      <div class="flex flex-col items-center">
        <p class="text-xl">6</p>
        <p class="text-center text-xs">In Progress</p>
      </div>
    </div>

    <div class="p-6">
      <div class="flex items-center justify-between">
        <h2 class="text-lg font-semibold">Continue</h2>
        <Button size="sm" variant="link" href="#">View More</Button>
      </div>

      <ScrollArea orientation="horizontal">
        <div class="flex w-full gap-2 py-4">
          {#each data.userMedia as media}
            <SmallMediaCard
              href="/media/{media.id}"
              coverUrl={media.coverUrl}
              title={media.title}
              startDate={undefined}
              partCount={0}
              score={undefined}
              userList={media.user?.list}
            />
          {/each}
        </div>
      </ScrollArea>
    </div>

    <div class="p-6">
      <div class="flex items-center justify-between">
        <h2 class="text-lg font-semibold">Recently created media</h2>
        <Button size="sm" variant="link" href="#">View More</Button>
      </div>

      <ScrollArea orientation="horizontal">
        <div class="flex w-full gap-2 py-4">
          {#each data.recentlyCreatedMedia as media}
            <SmallMediaCard
              href="/media/{media.id}"
              coverUrl={media.coverUrl}
              title={media.title}
              startDate={undefined}
              partCount={0}
              score={undefined}
              userList={media.user?.list}
            />
          {/each}
        </div>
      </ScrollArea>
    </div>

    <div class="p-6">
      <div class="flex items-center justify-between">
        <h2 class="text-lg font-semibold">Recently created collections</h2>
        <Button size="sm" variant="link" href="#">View More</Button>
      </div>

      <ScrollArea orientation="horizontal">
        <div class="flex w-full gap-2 py-4">
          {#each data.recentlyCreatedCollections as collection}
            <SmallMediaCard
              href="/collections/{collection.id}"
              coverUrl={collection.coverUrl}
              title={collection.name}
              startDate={undefined}
              partCount={0}
              score={undefined}
              userList={undefined}
            />
          {/each}
        </div>
      </ScrollArea>
    </div>
  </Card.Root>
{/if}

<div class="mx-auto max-w-5xl">
  <!-- Minimal Header -->
  <header class="mb-6 border-b border-gray-200 py-6 text-center">
    <h1 class="text-3xl font-bold text-gray-900">John</h1>
    <p class="text-sm text-gray-500">Private media library</p>
  </header>

  <!-- Main Content -->
  <div class="p-6">
    <!-- Unified Dashboard Card -->
    <div class="space-y-8 rounded-2xl bg-white p-6 shadow">
      <!-- Stats Row -->
      <section class="grid grid-cols-2 gap-6 text-center sm:grid-cols-4">
        <div>
          <p class="text-xl font-bold text-indigo-600">128</p>
          <p class="text-xs text-gray-600">Completed</p>
        </div>
        <div>
          <p class="text-xl font-bold text-indigo-600">6</p>
          <p class="text-xs text-gray-600">Watching</p>
        </div>
        <div>
          <p class="text-xl font-bold text-indigo-600">34</p>
          <p class="text-xs text-gray-600">Plan to Watch</p>
        </div>
        <div>
          <p class="text-xl font-bold text-indigo-600">412h</p>
          <p class="text-xs text-gray-600">Time Watched</p>
        </div>
      </section>

      <!-- Continue Watching -->
      <section>
        <h2 class="mb-3 text-lg font-semibold text-gray-900">
          Continue Watching
        </h2>
        <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4">
          <div class="overflow-hidden rounded-lg bg-gray-100">
            <div class="h-28 bg-gray-300"></div>
            <div class="p-2">
              <p class="text-sm font-medium text-gray-800">Attack on Titan</p>
              <p class="text-xs text-gray-500">Ep 14 of 25</p>
            </div>
          </div>
          <div class="overflow-hidden rounded-lg bg-gray-100">
            <div class="h-28 bg-gray-300"></div>
            <div class="p-2">
              <p class="text-sm font-medium text-gray-800">Steins;Gate</p>
              <p class="text-xs text-gray-500">Ep 8 of 24</p>
            </div>
          </div>
        </div>
      </section>

      <!-- Recent Updates -->
      <section>
        <h2 class="mb-3 text-lg font-semibold text-gray-900">
          Recent Updates
        </h2>
        <ul class="divide-y divide-gray-200">
          <li class="flex justify-between py-3">
            <span
              ><strong class="text-gray-800">Made in Abyss</strong> – Completed</span
            >
            <span class="text-sm text-gray-500">2d ago</span>
          </li>
          <li class="flex justify-between py-3">
            <span
              ><strong class="text-gray-800">Vinland Saga</strong> – Watching Ep
              6</span
            >
            <span class="text-sm text-gray-500">5d ago</span>
          </li>
          <li class="flex justify-between py-3">
            <span
              ><strong class="text-gray-800">One Piece</strong> – Added to Plan
              to Watch</span
            >
            <span class="text-sm text-gray-500">1w ago</span>
          </li>
        </ul>
      </section>

      <!-- User Lists -->
      <section>
        <h2 class="mb-3 text-lg font-semibold text-gray-900">Your Lists</h2>
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3">
          <div class="rounded-lg bg-gray-100 p-4">
            <h3 class="font-medium text-gray-800">Favorites</h3>
            <p class="text-xs text-gray-500">12 items</p>
          </div>
          <div class="rounded-lg bg-gray-100 p-4">
            <h3 class="font-medium text-gray-800">Rewatch Later</h3>
            <p class="text-xs text-gray-500">5 items</p>
          </div>
          <div class="rounded-lg bg-gray-100 p-4">
            <h3 class="font-medium text-gray-800">Top Sci-Fi</h3>
            <p class="text-xs text-gray-500">8 items</p>
          </div>
        </div>
      </section>
    </div>
  </div>
</div>

<div class="mx-auto max-w-5xl">
  <!-- Gradient Header -->
  <div
    class="flex h-32 items-end rounded-b-2xl bg-gradient-to-r from-indigo-500 via-purple-500 to-pink-500 p-6"
  >
    <div>
      <h1 class="text-3xl font-bold text-white">John</h1>
      <p class="text-sm text-gray-100">Private media library</p>
    </div>
  </div>

  <!-- Main Content -->
  <div class="p-6">
    <!-- Unified Dashboard Card -->
    <div class="space-y-8 rounded-2xl bg-white p-6 shadow">
      <!-- Stats Row -->
      <section class="grid grid-cols-2 gap-6 text-center sm:grid-cols-4">
        <div>
          <p class="text-xl font-bold text-indigo-600">128</p>
          <p class="text-xs text-gray-600">Completed</p>
        </div>
        <div>
          <p class="text-xl font-bold text-indigo-600">6</p>
          <p class="text-xs text-gray-600">Watching</p>
        </div>
        <div>
          <p class="text-xl font-bold text-indigo-600">34</p>
          <p class="text-xs text-gray-600">Plan to Watch</p>
        </div>
        <div>
          <p class="text-xl font-bold text-indigo-600">412h</p>
          <p class="text-xs text-gray-600">Time Watched</p>
        </div>
      </section>

      <!-- Continue Watching -->
      <section>
        <h2 class="mb-3 text-lg font-semibold text-gray-900">
          Continue Watching
        </h2>
        <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4">
          <div class="overflow-hidden rounded-lg bg-gray-100">
            <div class="h-28 bg-gray-300"></div>
            <div class="p-2">
              <p class="text-sm font-medium text-gray-800">Attack on Titan</p>
              <p class="text-xs text-gray-500">Ep 14 of 25</p>
            </div>
          </div>
          <div class="overflow-hidden rounded-lg bg-gray-100">
            <div class="h-28 bg-gray-300"></div>
            <div class="p-2">
              <p class="text-sm font-medium text-gray-800">Steins;Gate</p>
              <p class="text-xs text-gray-500">Ep 8 of 24</p>
            </div>
          </div>
        </div>
      </section>

      <!-- Recent Updates -->
      <section>
        <h2 class="mb-3 text-lg font-semibold text-gray-900">
          Recent Updates
        </h2>
        <ul class="divide-y divide-gray-200">
          <li class="flex justify-between py-3">
            <span
              ><strong class="text-gray-800">Made in Abyss</strong> – Completed</span
            >
            <span class="text-sm text-gray-500">2d ago</span>
          </li>
          <li class="flex justify-between py-3">
            <span
              ><strong class="text-gray-800">Vinland Saga</strong> – Watching Ep
              6</span
            >
            <span class="text-sm text-gray-500">5d ago</span>
          </li>
          <li class="flex justify-between py-3">
            <span
              ><strong class="text-gray-800">One Piece</strong> – Added to Plan
              to Watch</span
            >
            <span class="text-sm text-gray-500">1w ago</span>
          </li>
        </ul>
      </section>

      <!-- User Lists -->
      <section>
        <h2 class="mb-3 text-lg font-semibold text-gray-900">Your Lists</h2>
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3">
          <div class="rounded-lg bg-gray-100 p-4">
            <h3 class="font-medium text-gray-800">Favorites</h3>
            <p class="text-xs text-gray-500">12 items</p>
          </div>
          <div class="rounded-lg bg-gray-100 p-4">
            <h3 class="font-medium text-gray-800">Rewatch Later</h3>
            <p class="text-xs text-gray-500">5 items</p>
          </div>
          <div class="rounded-lg bg-gray-100 p-4">
            <h3 class="font-medium text-gray-800">Top Sci-Fi</h3>
            <p class="text-xs text-gray-500">8 items</p>
          </div>
        </div>
      </section>
    </div>
  </div>
</div>

<div class="mx-auto max-w-5xl p-6">
  <!-- Page Header -->
  <header class="mb-6">
    <h1 class="text-3xl font-bold text-gray-900">Welcome back, John</h1>
    <p class="text-sm text-gray-600">Your personal watch dashboard</p>
  </header>

  <!-- Unified Dashboard Card -->
  <div class="space-y-8 rounded-2xl bg-white p-6 shadow">
    <!-- Stats Row -->
    <section class="grid grid-cols-2 gap-6 text-center sm:grid-cols-4">
      <div>
        <p class="text-xl font-bold text-indigo-600">128</p>
        <p class="text-xs text-gray-600">Completed</p>
      </div>
      <div>
        <p class="text-xl font-bold text-indigo-600">6</p>
        <p class="text-xs text-gray-600">Watching</p>
      </div>
      <div>
        <p class="text-xl font-bold text-indigo-600">34</p>
        <p class="text-xs text-gray-600">Plan to Watch</p>
      </div>
      <div>
        <p class="text-xl font-bold text-indigo-600">412h</p>
        <p class="text-xs text-gray-600">Time Watched</p>
      </div>
    </section>

    <!-- Continue Watching -->
    <section>
      <h2 class="mb-3 text-lg font-semibold text-gray-900">
        Continue Watching
      </h2>
      <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4">
        <div class="overflow-hidden rounded-lg bg-gray-100">
          <div class="h-28 bg-gray-300"></div>
          <div class="p-2">
            <p class="text-sm font-medium text-gray-800">Attack on Titan</p>
            <p class="text-xs text-gray-500">Ep 14 of 25</p>
          </div>
        </div>
        <div class="overflow-hidden rounded-lg bg-gray-100">
          <div class="h-28 bg-gray-300"></div>
          <div class="p-2">
            <p class="text-sm font-medium text-gray-800">Steins;Gate</p>
            <p class="text-xs text-gray-500">Ep 8 of 24</p>
          </div>
        </div>
        <!-- more items -->
      </div>
    </section>

    <!-- Recent Updates -->
    <section>
      <h2 class="mb-3 text-lg font-semibold text-gray-900">Recent Updates</h2>
      <ul class="divide-y divide-gray-200">
        <li class="flex justify-between py-3">
          <span
            ><strong class="text-gray-800">Made in Abyss</strong> – Completed</span
          >
          <span class="text-sm text-gray-500">2d ago</span>
        </li>
        <li class="flex justify-between py-3">
          <span
            ><strong class="text-gray-800">Vinland Saga</strong> – Watching Ep 6</span
          >
          <span class="text-sm text-gray-500">5d ago</span>
        </li>
        <li class="flex justify-between py-3">
          <span
            ><strong class="text-gray-800">One Piece</strong> – Added to Plan to
            Watch</span
          >
          <span class="text-sm text-gray-500">1w ago</span>
        </li>
      </ul>
    </section>

    <!-- User Lists -->
    <section>
      <h2 class="mb-3 text-lg font-semibold text-gray-900">Your Lists</h2>
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3">
        <div class="rounded-lg bg-gray-100 p-4">
          <h3 class="font-medium text-gray-800">Favorites</h3>
          <p class="text-xs text-gray-500">12 items</p>
        </div>
        <div class="rounded-lg bg-gray-100 p-4">
          <h3 class="font-medium text-gray-800">Rewatch Later</h3>
          <p class="text-xs text-gray-500">5 items</p>
        </div>
        <div class="rounded-lg bg-gray-100 p-4">
          <h3 class="font-medium text-gray-800">Top Sci-Fi</h3>
          <p class="text-xs text-gray-500">8 items</p>
        </div>
      </div>
    </section>
  </div>
</div>

<div class="h-[400px]"></div>

<div
  class="flex min-h-screen flex-col items-center justify-center bg-gray-50 px-6"
>
  <!-- App Title -->
  <h1 class="mb-4 text-4xl font-bold text-indigo-600">WatchBook</h1>
  <p class="mb-8 max-w-md text-center text-gray-600">
    Keep track of your shows and movies in your own private library.
  </p>

  <!-- Auth Buttons -->
  <div class="space-x-4">
    <a
      href="/login"
      class="rounded-lg bg-indigo-600 px-5 py-2 text-white shadow hover:bg-indigo-700"
    >
      Login
    </a>
    <a
      href="/register"
      class="rounded-lg border border-gray-300 px-5 py-2 text-gray-700 hover:bg-gray-100"
    >
      Register
    </a>
  </div>
</div>

<div class="min-h-screen bg-gray-50">
  <!-- Header -->
  <header class="mx-auto flex max-w-7xl items-center justify-between p-6">
    <h1 class="text-2xl font-bold text-indigo-600">WatchBook</h1>
    <div class="space-x-4">
      <a href="/login" class="text-sm text-gray-700 hover:text-indigo-600"
        >Login</a
      >
      <a
        href="/register"
        class="rounded-lg bg-indigo-600 px-4 py-2 text-sm text-white shadow hover:bg-indigo-700"
      >
        Register
      </a>
    </div>
  </header>

  <!-- Hero Section -->
  <section class="mx-auto max-w-4xl px-6 py-20 text-center">
    <h2 class="mb-4 text-4xl font-bold text-gray-900">
      Track and Organize Your Watchlist
    </h2>
    <p class="mb-6 text-lg text-gray-600">
      A self-hosted app to keep track of what you're watching. Create custom
      folders, track progress, and manage your own private media library.
    </p>
    <a
      href="/register"
      class="rounded-xl bg-indigo-600 px-6 py-3 text-lg text-white shadow hover:bg-indigo-700"
    >
      Get Started
    </a>
  </section>

  <!-- Features Section -->
  <section
    class="mx-auto grid max-w-6xl gap-8 px-6 py-16 sm:grid-cols-2 lg:grid-cols-3"
  >
    <div class="rounded-2xl bg-white p-6 shadow">
      <h3 class="mb-2 text-lg font-semibold text-gray-900">Track Progress</h3>
      <p class="text-sm text-gray-600">
        Easily update what you're watching, completed, or plan to watch.
      </p>
    </div>
    <div class="rounded-2xl bg-white p-6 shadow">
      <h3 class="mb-2 text-lg font-semibold text-gray-900">Custom Folders</h3>
      <p class="text-sm text-gray-600">
        Organize your favorites with flexible, user-created folders.
      </p>
    </div>
    <div class="rounded-2xl bg-white p-6 shadow">
      <h3 class="mb-2 text-lg font-semibold text-gray-900">
        Private & Self-Hosted
      </h3>
      <p class="text-sm text-gray-600">
        All data is yours — host your own watch tracker without relying on
        external services.
      </p>
    </div>
  </section>

  <!-- Preview Section -->
  <section class="bg-indigo-50 py-16">
    <div class="mx-auto max-w-6xl px-6 text-center">
      <h3 class="mb-6 text-2xl font-bold text-gray-900">See Your Dashboard</h3>
      <div
        class="flex h-64 items-center justify-center rounded-2xl bg-white text-gray-400 shadow-lg"
      >
        [Preview Screenshot Placeholder]
      </div>
    </div>
  </section>

  <!-- Footer -->
  <footer
    class="mx-auto max-w-7xl px-6 py-10 text-center text-sm text-gray-500"
  >
    © 2025 WatchBook. Built for self-hosted media tracking.
  </footer>
</div>

<div class="mx-auto max-w-7xl p-6">
  <!-- Header -->
  <header class="mb-6">
    <h1 class="text-3xl font-bold text-gray-900">Welcome back, John</h1>
    <p class="text-sm text-gray-600">Your personal watch dashboard</p>
  </header>

  <!-- Main Layout -->
  <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
    <!-- Left Sidebar (Lists & Folders) -->
    <aside class="space-y-6 lg:col-span-1">
      <!-- System Lists -->
      <div class="rounded-2xl bg-white p-4 shadow">
        <h2 class="mb-3 text-lg font-semibold text-gray-900">My Lists</h2>
        <ul class="space-y-2 text-sm text-gray-700">
          <li class="flex justify-between">
            <span>Watching</span>
            <span class="text-gray-500">6</span>
          </li>
          <li class="flex justify-between">
            <span>Completed</span>
            <span class="text-gray-500">128</span>
          </li>
          <li class="flex justify-between">
            <span>Plan to Watch</span>
            <span class="text-gray-500">34</span>
          </li>
          <li class="flex justify-between">
            <span>Dropped</span>
            <span class="text-gray-500">3</span>
          </li>
        </ul>
      </div>

      <!-- User Folders -->
      <div class="rounded-2xl bg-white p-4 shadow">
        <h2 class="mb-3 text-lg font-semibold text-gray-900">Folders</h2>
        <ul class="space-y-2 text-sm text-gray-700">
          <li>Favorites (12)</li>
          <li>Rewatch Later (5)</li>
          <li>Top Sci-Fi (8)</li>
        </ul>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="space-y-8 lg:col-span-2">
      <!-- Stats -->
      <section class="grid grid-cols-2 gap-4 sm:grid-cols-4">
        <div class="rounded-2xl bg-white p-4 text-center shadow">
          <p class="text-xl font-bold text-indigo-600">128</p>
          <p class="text-xs text-gray-600">Completed</p>
        </div>
        <div class="rounded-2xl bg-white p-4 text-center shadow">
          <p class="text-xl font-bold text-indigo-600">6</p>
          <p class="text-xs text-gray-600">Watching</p>
        </div>
        <div class="rounded-2xl bg-white p-4 text-center shadow">
          <p class="text-xl font-bold text-indigo-600">34</p>
          <p class="text-xs text-gray-600">Plan to Watch</p>
        </div>
        <div class="rounded-2xl bg-white p-4 text-center shadow">
          <p class="text-xl font-bold text-indigo-600">412h</p>
          <p class="text-xs text-gray-600">Time Watched</p>
        </div>
      </section>

      <!-- Continue Watching -->
      <section>
        <h2 class="mb-3 text-lg font-semibold text-gray-900">
          Continue Watching
        </h2>
        <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
          <div class="overflow-hidden rounded-lg bg-white shadow">
            <div class="h-32 bg-gray-200"></div>
            <div class="p-3">
              <p class="text-sm font-medium text-gray-800">Attack on Titan</p>
              <p class="text-xs text-gray-500">Ep 14 of 25</p>
            </div>
          </div>
          <div class="overflow-hidden rounded-lg bg-white shadow">
            <div class="h-32 bg-gray-200"></div>
            <div class="p-3">
              <p class="text-sm font-medium text-gray-800">Steins;Gate</p>
              <p class="text-xs text-gray-500">Ep 8 of 24</p>
            </div>
          </div>
          <!-- more cards -->
        </div>
      </section>

      <!-- Recent Updates -->
      <section>
        <h2 class="mb-3 text-lg font-semibold text-gray-900">
          Recent Updates
        </h2>
        <ul class="divide-y divide-gray-200 rounded-2xl bg-white shadow">
          <li class="flex justify-between p-4">
            <span
              ><strong class="text-gray-800">Made in Abyss</strong> – Completed</span
            >
            <span class="text-sm text-gray-500">2d ago</span>
          </li>
          <li class="flex justify-between p-4">
            <span
              ><strong class="text-gray-800">Vinland Saga</strong> – Watching Ep
              6</span
            >
            <span class="text-sm text-gray-500">5d ago</span>
          </li>
          <li class="flex justify-between p-4">
            <span
              ><strong class="text-gray-800">One Piece</strong> – Added to Plan
              to Watch</span
            >
            <span class="text-sm text-gray-500">1w ago</span>
          </li>
        </ul>
      </section>
    </main>
  </div>
</div>

<div class="mx-auto max-w-7xl space-y-8 p-6">
  <!-- Header -->
  <header>
    <h1 class="text-3xl font-bold text-gray-900">Welcome back, John</h1>
    <p class="text-sm text-gray-600">Here’s what’s new in your watchlist.</p>
  </header>

  <!-- Quick Stats -->
  <section class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
    <div class="rounded-2xl bg-white p-6 text-center shadow">
      <p class="text-2xl font-bold text-indigo-600">128</p>
      <p class="text-sm text-gray-600">Completed</p>
    </div>
    <div class="rounded-2xl bg-white p-6 text-center shadow">
      <p class="text-2xl font-bold text-indigo-600">6</p>
      <p class="text-sm text-gray-600">Currently Watching</p>
    </div>
    <div class="rounded-2xl bg-white p-6 text-center shadow">
      <p class="text-2xl font-bold text-indigo-600">34</p>
      <p class="text-sm text-gray-600">Plan to Watch</p>
    </div>
    <div class="rounded-2xl bg-white p-6 text-center shadow">
      <p class="text-2xl font-bold text-indigo-600">412h</p>
      <p class="text-sm text-gray-600">Total Time Watched</p>
    </div>
  </section>

  <!-- Continue Watching -->
  <section>
    <h2 class="mb-3 text-lg font-semibold text-gray-900">Continue Watching</h2>
    <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
      <div class="overflow-hidden rounded-lg bg-white shadow">
        <div class="h-32 bg-gray-200"></div>
        <div class="p-3">
          <p class="text-sm font-medium text-gray-800">Attack on Titan</p>
          <p class="text-xs text-gray-500">Ep 14 of 25</p>
        </div>
      </div>
      <div class="overflow-hidden rounded-lg bg-white shadow">
        <div class="h-32 bg-gray-200"></div>
        <div class="p-3">
          <p class="text-sm font-medium text-gray-800">Steins;Gate</p>
          <p class="text-xs text-gray-500">Ep 8 of 24</p>
        </div>
      </div>
      <!-- more cards -->
    </div>
  </section>

  <!-- Recent Updates -->
  <section>
    <h2 class="mb-3 text-lg font-semibold text-gray-900">Recent Updates</h2>
    <ul class="divide-y divide-gray-200 rounded-2xl bg-white shadow">
      <li class="flex justify-between p-4">
        <span
          ><strong class="text-gray-800">Made in Abyss</strong> – Completed</span
        >
        <span class="text-sm text-gray-500">2d ago</span>
      </li>
      <li class="flex justify-between p-4">
        <span
          ><strong class="text-gray-800">Vinland Saga</strong> – Watching Ep 6</span
        >
        <span class="text-sm text-gray-500">5d ago</span>
      </li>
      <li class="flex justify-between p-4">
        <span
          ><strong class="text-gray-800">One Piece</strong> – Added to Plan to Watch</span
        >
        <span class="text-sm text-gray-500">1w ago</span>
      </li>
    </ul>
  </section>

  <!-- Favorite Lists -->
  <section>
    <h2 class="mb-3 text-lg font-semibold text-gray-900">Your Lists</h2>
    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <div class="rounded-xl bg-white p-4 shadow">
        <h3 class="font-medium text-gray-800">Favorites</h3>
        <p class="text-xs text-gray-500">12 items</p>
      </div>
      <div class="rounded-xl bg-white p-4 shadow">
        <h3 class="font-medium text-gray-800">Rewatch Later</h3>
        <p class="text-xs text-gray-500">5 items</p>
      </div>
      <div class="rounded-xl bg-white p-4 shadow">
        <h3 class="font-medium text-gray-800">Top Sci-Fi</h3>
        <p class="text-xs text-gray-500">8 items</p>
      </div>
    </div>
  </section>
</div>

<div class="space-y-8 p-4">
  <!-- Hero / Featured Banner -->
  <section class="relative overflow-hidden rounded-xl shadow-lg">
    <img
      src="https://via.placeholder.com/1200x400"
      alt="Featured Banner"
      class="h-[250px] w-full object-cover"
    />
    <div
      class="absolute inset-0 flex flex-col justify-end bg-gradient-to-t from-black/80 to-transparent p-6"
    >
      <h1 class="text-2xl font-bold text-white">
        Top Pick: Attack on Titan Final Season
      </h1>
      <p class="mt-1 text-sm text-gray-300">
        The epic conclusion — streaming now.
      </p>
    </div>
  </section>

  <!-- Featured Collections -->
  <section>
    <h2 class="mb-3 text-xl font-bold">Collections</h2>
    <div
      class="scrollbar-thin scrollbar-thumb-gray-600 flex gap-4 overflow-x-auto pb-2"
    >
      <a
        href="#"
        class="w-64 flex-shrink-0 overflow-hidden rounded-lg bg-gray-800 shadow transition hover:shadow-lg"
      >
        <img
          src="https://via.placeholder.com/256x144"
          alt="Collection Cover"
          class="h-36 w-full object-cover"
        />
        <div class="p-4">
          <h3 class="truncate text-lg font-semibold">Collection Name</h3>
        </div>
      </a>
      <!-- Repeat -->
    </div>
  </section>

  <!-- Trending Media -->
  <section>
    <h2 class="mb-3 text-xl font-bold">Trending</h2>
    <div
      class="scrollbar-thin scrollbar-thumb-gray-600 flex gap-4 overflow-x-auto pb-2"
    >
      <a
        href="#"
        class="w-40 flex-shrink-0 overflow-hidden rounded-lg bg-gray-800 shadow transition hover:shadow-lg"
      >
        <img
          src="https://via.placeholder.com/160x240"
          alt="Media Cover"
          class="h-60 w-full object-cover"
        />
        <div class="p-3">
          <h3 class="line-clamp-2 text-sm font-semibold">Anime Title</h3>
          <div class="mt-1 flex gap-2 text-xs text-gray-400">
            <span>12 eps</span>
            <span>2023</span>
          </div>
        </div>
      </a>
      <!-- Repeat -->
    </div>
  </section>

  <!-- Your Activity -->
  <section>
    <h2 class="mb-3 text-xl font-bold">Your Activity</h2>
    <div class="grid grid-cols-1 gap-4 md:grid-cols-3">
      <!-- Currently Watching -->
      <div class="rounded-lg bg-gray-800 p-4 shadow">
        <h3 class="mb-3 text-lg font-semibold">Currently Watching</h3>
        <ul class="space-y-2">
          <li class="flex items-center gap-3">
            <img
              src="https://via.placeholder.com/60x90"
              class="h-[90px] w-[60px] rounded object-cover"
              alt=""
            />
            <div>
              <p class="font-medium">Anime Title</p>
              <p class="text-sm text-gray-400">Ep 5 / 12</p>
            </div>
          </li>
        </ul>
      </div>

      <!-- Backlog -->
      <div class="rounded-lg bg-gray-800 p-4 shadow">
        <h3 class="mb-3 text-lg font-semibold">Backlog</h3>
        <ul class="space-y-2">
          <li class="flex items-center gap-3">
            <img
              src="https://via.placeholder.com/60x90"
              class="h-[90px] w-[60px] rounded object-cover"
              alt=""
            />
            <div>
              <p class="font-medium">Anime Title</p>
              <p class="text-sm text-gray-400">Plan to Watch</p>
            </div>
          </li>
        </ul>
      </div>

      <!-- Recently Completed -->
      <div class="rounded-lg bg-gray-800 p-4 shadow">
        <h3 class="mb-3 text-lg font-semibold">Recently Completed</h3>
        <ul class="space-y-2">
          <li class="flex items-center gap-3">
            <img
              src="https://via.placeholder.com/60x90"
              class="h-[90px] w-[60px] rounded object-cover"
              alt=""
            />
            <div>
              <p class="font-medium">Anime Title</p>
              <p class="text-sm text-gray-400">Completed</p>
            </div>
          </li>
        </ul>
      </div>
    </div>
  </section>
</div>

<div class="space-y-8 p-4">
  <!-- Featured Collections -->
  <section>
    <h2 class="mb-4 text-xl font-bold">Featured Collections</h2>
    <div
      class="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4"
    >
      <!-- Example Collection Card -->
      <a
        href="#"
        class="w-full max-w-[320px] overflow-hidden rounded-lg bg-gray-800 shadow transition hover:shadow-lg"
      >
        <img
          src="https://via.placeholder.com/320x180"
          alt="Collection Cover"
          class="h-48 w-full object-cover"
        />
        <div class="p-4">
          <h3 class="truncate text-lg font-semibold">Collection Name</h3>
          <span
            class="mt-2 inline-block rounded-full bg-blue-500 px-2 py-1 text-sm text-white"
            >Anime</span
          >
        </div>
      </a>
      <!-- Repeat -->
    </div>
  </section>

  <!-- Featured Media Items -->
  <section>
    <h2 class="mb-4 text-xl font-bold">Featured Media</h2>
    <div
      class="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5"
    >
      <!-- Example Media Card -->
      <a
        href="#"
        class="w-full max-w-[200px] overflow-hidden rounded-lg bg-gray-800 shadow transition hover:shadow-lg"
      >
        <img
          src="https://via.placeholder.com/200x300"
          alt="Media Cover"
          class="h-[300px] w-full object-cover"
        />
        <div class="p-4">
          <h3 class="line-clamp-2 text-sm font-semibold">
            Anime Season Title Goes Here
          </h3>
          <div class="mt-2 flex gap-2 text-xs text-gray-400">
            <span>12 eps</span>
            <span>2023</span>
          </div>
        </div>
      </a>
      <!-- Repeat -->
    </div>
  </section>

  <!-- User Data -->
  <section>
    <h2 class="mb-4 text-xl font-bold">Your Activity</h2>
    <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
      <!-- Currently Watching -->
      <div class="rounded-lg bg-gray-800 p-4 shadow">
        <h3 class="mb-3 text-lg font-semibold">Currently Watching</h3>
        <ul class="space-y-2">
          <li class="flex items-center gap-3">
            <img
              src="https://via.placeholder.com/60x90"
              class="h-[90px] w-[60px] rounded object-cover"
              alt=""
            />
            <div>
              <p class="font-medium">Anime Title</p>
              <p class="text-sm text-gray-400">Episode 5 / 12</p>
            </div>
          </li>
        </ul>
      </div>

      <!-- Backlog -->
      <div class="rounded-lg bg-gray-800 p-4 shadow">
        <h3 class="mb-3 text-lg font-semibold">Backlog</h3>
        <ul class="space-y-2">
          <li class="flex items-center gap-3">
            <img
              src="https://via.placeholder.com/60x90"
              class="h-[90px] w-[60px] rounded object-cover"
              alt=""
            />
            <div>
              <p class="font-medium">Anime Title</p>
              <p class="text-sm text-gray-400">Plan to Watch</p>
            </div>
          </li>
        </ul>
      </div>
    </div>
  </section>
</div>

<div class="space-y-8 p-4">
  <!-- Search + Filters -->
  <div class="rounded-lg bg-gray-800 p-4 shadow-lg">
    <div
      class="flex flex-col space-y-3 md:flex-row md:items-center md:space-x-4 md:space-y-0"
    >
      <input
        type="text"
        placeholder="Search..."
        class="flex-1 rounded-lg bg-gray-700 px-4 py-2 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
      <select
        class="rounded-lg bg-gray-700 px-4 py-2 text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="">All Types</option>
        <option value="anime">Anime</option>
        <option value="tv">TV</option>
        <option value="game">Game</option>
        <option value="more">More</option>
      </select>
      <select
        class="rounded-lg bg-gray-700 px-4 py-2 text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="latest">Latest</option>
        <option value="oldest">Oldest</option>
        <option value="az">A-Z</option>
        <option value="za">Z-A</option>
      </select>
    </div>
  </div>

  <!-- Featured Collections -->
  <section>
    <h2 class="mb-3 text-xl font-bold text-white">Featured Collections</h2>
    <div
      class="grid justify-items-center gap-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4"
    >
      <!-- Card -->
      <a
        href="#"
        class="w-full max-w-[320px] overflow-hidden rounded-lg bg-gray-900 shadow-lg transition-transform hover:scale-[1.02]"
      >
        <div class="relative">
          <img
            src="https://via.placeholder.com/320x180"
            alt="Cover"
            class="h-48 w-full object-cover"
          />
          <div
            class="absolute bottom-0 w-full bg-black/50 p-2 backdrop-blur-md"
          >
            <span class="block truncate font-semibold text-white"
              >Collection Title</span
            >
            <span class="text-xs text-blue-400">Plan to Watch</span>
          </div>
        </div>
      </a>
      <!-- Repeat -->
    </div>
  </section>

  <!-- Recent Media -->
  <section>
    <h2 class="mb-3 text-xl font-bold text-white">Recently Added Media</h2>
    <div
      class="grid justify-items-center gap-6 sm:grid-cols-2 md:grid-cols-4 lg:grid-cols-6"
    >
      <a
        href="#"
        class="w-full max-w-[180px] overflow-hidden rounded-lg bg-gray-900 shadow-lg transition-transform hover:scale-[1.02]"
      >
        <img
          src="https://via.placeholder.com/180x260"
          alt="Media"
          class="h-[260px] w-full object-cover"
        />
        <div class="p-2">
          <span class="block truncate text-sm font-semibold text-white"
            >Media Title</span
          >
          <span class="text-xs text-gray-400">Anime</span>
        </div>
      </a>
    </div>
  </section>

  <!-- User Dashboard -->
  <section>
    <h2 class="mb-3 text-xl font-bold text-white">Your Watchlist</h2>
    <div
      class="grid justify-items-center gap-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4"
    >
      <a
        href="#"
        class="w-full max-w-[320px] overflow-hidden rounded-lg bg-gray-900 shadow-lg transition-transform hover:scale-[1.02]"
      >
        <div class="relative">
          <img
            src="https://via.placeholder.com/320x180"
            alt="Cover"
            class="h-48 w-full object-cover"
          />
          <div
            class="absolute bottom-0 w-full bg-black/50 p-2 backdrop-blur-md"
          >
            <span class="block truncate font-semibold text-white"
              >Attack on Titan</span
            >
            <span class="text-xs text-green-400"
              >Currently Watching - Ep 15/25</span
            >
          </div>
        </div>
      </a>
      <a
        href="#"
        class="w-full max-w-[320px] overflow-hidden rounded-lg bg-gray-900 shadow-lg transition-transform hover:scale-[1.02]"
      >
        <div class="relative">
          <img
            src="https://via.placeholder.com/320x180"
            alt="Cover"
            class="h-48 w-full object-cover"
          />
          <div
            class="absolute bottom-0 w-full bg-black/50 p-2 backdrop-blur-md"
          >
            <span class="block truncate font-semibold text-white"
              >One Piece</span
            >
            <span class="text-xs text-yellow-400">Backlog - Ep 230/1070</span>
          </div>
        </div>
      </a>
    </div>
  </section>
</div>

<div class="space-y-6 p-4">
  <!-- Search and Filters -->
  <div class="rounded-lg bg-gray-800 p-4 shadow-lg">
    <div
      class="flex flex-col space-y-3 md:flex-row md:items-center md:space-x-4 md:space-y-0"
    >
      <!-- Search -->
      <input
        type="text"
        placeholder="Search collections..."
        class="flex-1 rounded-lg bg-gray-700 px-4 py-2 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
      />

      <!-- Type Filter -->
      <select
        class="rounded-lg bg-gray-700 px-4 py-2 text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="">All Types</option>
        <option value="anime">Anime</option>
        <option value="tv">TV</option>
        <option value="game">Game</option>
        <option value="more">More</option>
      </select>

      <!-- Sort -->
      <select
        class="rounded-lg bg-gray-700 px-4 py-2 text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="latest">Latest</option>
        <option value="oldest">Oldest</option>
        <option value="az">A-Z</option>
        <option value="za">Z-A</option>
      </select>
    </div>
  </div>

  <!-- Results Count -->
  <div class="text-sm text-gray-400">Showing 24 collections</div>

  <!-- Grid -->
  <div
    class="grid justify-items-center gap-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4"
  >
    <!-- Card -->
    <a
      href="#"
      class="w-full max-w-[320px] overflow-hidden rounded-lg bg-gray-900 shadow-lg transition-transform hover:scale-[1.02]"
    >
      <div class="relative">
        <img
          src="https://via.placeholder.com/320x180"
          alt="Cover"
          class="h-48 w-full object-cover"
        />
        <div class="absolute bottom-0 w-full bg-black/50 p-2 backdrop-blur-md">
          <span class="block truncate font-semibold text-white"
            >Collection Title</span
          >
          <span class="text-xs text-blue-400">Plan to Watch</span>
        </div>
      </div>
    </a>

    <!-- Repeat more cards -->
  </div>
</div>
