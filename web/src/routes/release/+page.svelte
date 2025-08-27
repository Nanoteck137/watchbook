<script>
  import Image from "$lib/components/Image.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
</script>

<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
  <div
    class="grid grid-cols-[repeat(auto-fit,minmax(240px,1fr))] items-center justify-items-center gap-6"
  >
    <!-- svelte-ignore a11y_missing_attribute -->
    <!-- svelte-ignore a11y_invalid_attribute -->
    <a
      class="group relative aspect-[75/106] max-w-[240px] transform cursor-pointer overflow-hidden rounded-3xl bg-gray-800 shadow-md transition-transform duration-300 hover:-translate-y-1 hover:shadow-lg"
      href="#"
    >
      <!-- Badge -->
      <!-- {#if userList}
        {@const list = parseUserList(userList)}
        {#if list}
          <Badge class="absolute left-2 top-2 z-10" {list} />
        {/if}
      {/if} -->

      <div
        class="flex min-w-[240px] items-center justify-center bg-gray-800 text-gray-400"
      >
        <Image
          src="http://localhost:3000/files/collections/ohtm567u/images/cover.jpeg"
          alt="Cover image"
          class="aspect-[75/106] h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
        />
      </div>

      <div
        class="absolute bottom-0 left-0 right-0 rounded-t-lg bg-black/80 p-3 text-center backdrop-blur-md"
      >
        <h2
          class="line-clamp-2 text-ellipsis text-base font-semibold"
          title="Alya Sometimes Hides Her Feelings in Russian"
        >
          Alya Sometimes Hides Her Feelings in Russian
        </h2>

        <Spacer size="xs" />

        <div class="flex flex-col gap-2">
          <div class="flex flex-col gap-1">
            <p class="text-xs text-gray-400">5 / 12 eps</p>

            <div class="h-1 w-full rounded-full bg-gray-700">
              <div
                class="h-1 rounded-full bg-blue-500"
                style="width: 42%;"
              ></div>
            </div>
          </div>

          <div class="flex items-center gap-1">
            <p class="text-xs text-gray-300">Next in:</p>
            <p class="text-sm font-bold text-blue-400">3d 12h</p>
          </div>
        </div>
      </div>
    </a>

    <div
      class="flex w-36 flex-col overflow-hidden rounded-lg bg-gray-800 shadow transition-shadow hover:shadow-lg"
    >
      <div class="aspect-[75/106] w-full overflow-hidden">
        <img
          src="http://localhost:3000/files/collections/ohtm567u/images/cover.jpeg"
          alt="Show Cover"
          class="h-full w-full object-cover"
        />
      </div>
      <div class="flex flex-1 flex-col justify-between p-2 text-sm">
        <div>
          <h2 class="truncate text-base font-semibold leading-tight">
            My Show
          </h2>
          <p class="text-xs text-gray-400">5 / 12 eps</p>
          <!-- Progress bar -->
          <div class="mt-1 h-1 w-full rounded-full bg-gray-700">
            <div
              class="h-1 rounded-full bg-blue-500"
              style="width: 42%;"
            ></div>
          </div>
        </div>
        <div class="mt-2">
          <p class="text-xs text-gray-300">Next in:</p>
          <p class="text-sm font-bold text-blue-500">3d 12h</p>
        </div>
      </div>
    </div>
  </div>
</div>

<div class="mx-auto min-h-screen max-w-6xl bg-gray-900 p-4 text-gray-100">
  <h1 class="mb-6 text-center text-2xl font-bold">Release Calendar</h1>

  <div
    class="grid grid-cols-2 justify-items-center gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5"
  >
    <!-- Show Card -->
    <div
      class="flex w-36 flex-col overflow-hidden rounded-lg bg-gray-800 shadow transition-shadow hover:shadow-lg"
    >
      <div class="aspect-[75/106] w-full overflow-hidden">
        <img
          src="http://localhost:3000/files/collections/ohtm567u/images/cover.jpeg"
          alt="Show Cover"
          class="h-full w-full object-cover"
        />
      </div>
      <div class="flex flex-1 flex-col justify-between p-2 text-sm">
        <div>
          <h2 class="truncate text-base font-semibold leading-tight">
            My Show
          </h2>
          <p class="text-xs text-gray-400">5 / 12 eps</p>
          <!-- Progress bar -->
          <div class="mt-1 h-1 w-full rounded-full bg-gray-700">
            <div
              class="h-1 rounded-full bg-blue-500"
              style="width: 42%;"
            ></div>
          </div>
        </div>
        <div class="mt-2">
          <p class="text-xs text-gray-300">Next in:</p>
          <p class="text-sm font-bold text-blue-400">3d 12h</p>
        </div>
      </div>
    </div>

    <!-- Another Show -->
    <div
      class="flex w-36 flex-col overflow-hidden rounded-lg bg-gray-800 shadow transition-shadow hover:shadow-lg"
    >
      <div class="aspect-[75/106] w-full overflow-hidden">
        <img
          src="http://localhost:3000/files/collections/ohtm567u/images/cover.jpeg"
          alt="Show Cover"
          class="h-full w-full object-cover"
        />
      </div>
      <div class="flex flex-1 flex-col justify-between p-2 text-sm">
        <div>
          <h2 class="truncate text-base font-semibold leading-tight">
            Another Show
          </h2>
          <p class="text-xs text-gray-400">12 / 12 eps</p>
          <div class="mt-1 h-1 w-full rounded-full bg-gray-700">
            <div
              class="h-1 rounded-full bg-green-500"
              style="width: 100%;"
            ></div>
          </div>
        </div>
        <div class="mt-2">
          <p class="text-sm font-bold text-green-400">Completed</p>
        </div>
      </div>
    </div>

    <!-- Upcoming Show -->
    <div
      class="flex w-36 flex-col overflow-hidden rounded-lg bg-gray-800 shadow transition-shadow hover:shadow-lg"
    >
      <div class="aspect-[75/106] w-full overflow-hidden">
        <img
          src="http://localhost:3000/files/collections/ohtm567u/images/cover.jpeg"
          alt="Show Cover"
          class="h-full w-full object-cover"
        />
      </div>
      <div class="flex flex-1 flex-col justify-between p-2 text-sm">
        <div>
          <h2 class="truncate text-base font-semibold leading-tight">
            Upcoming Show
          </h2>
          <p class="text-xs text-gray-400">0 / 10 eps</p>
          <div class="mt-1 h-1 w-full rounded-full bg-gray-700">
            <div class="h-1 rounded-full bg-gray-500" style="width: 0%;"></div>
          </div>
        </div>
        <div class="mt-2">
          <p class="text-xs text-yellow-400">Starts in:</p>
          <p class="text-sm font-bold text-yellow-400">7d 4h</p>
        </div>
      </div>
    </div>
  </div>
</div>
