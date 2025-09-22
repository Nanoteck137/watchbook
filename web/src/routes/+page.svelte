<script>
  import { Button, Card, ScrollArea } from "@nanoteck137/nano-ui";
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
        class="h-full w-full bg-gradient-to-br from-[#9333EA] to-[#F59E0B] p-6"
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
        <Button
          size="sm"
          variant="link"
          href="/users/{data.user.id}/watchlist?list=in-progress"
          >View More</Button
        >
      </div>

      <ScrollArea orientation="horizontal">
        <div class="flex w-full gap-2 py-4">
          {#each data.userInprogressMedia as media}
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
        <h2 class="text-lg font-semibold">Backlog</h2>
        <Button
          size="sm"
          variant="link"
          href="/users/{data.user.id}/watchlist?list=backlog"
        >
          View More
        </Button>
      </div>

      <ScrollArea orientation="horizontal">
        <div class="flex w-full gap-2 py-4">
          {#each data.userBacklogMedia as media}
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
        <Button size="sm" variant="link" href="/media?sort=created-new">
          View More
        </Button>
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
        <Button size="sm" variant="link" href="/collections?sort=created-new">
          View More
        </Button>
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
