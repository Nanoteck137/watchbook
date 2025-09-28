<script lang="ts">
  import { Button, Card } from "@nanoteck137/nano-ui";
  import NewFolderModal from "./NewFolderModal.svelte";
  import { Plus } from "lucide-svelte";

  const { data } = $props();

  let openNewModal = $state(false);
</script>

<Card.Root>
  <div class="flex items-center gap-2 border-b p-6">
    <h1 class="text-xl">Folders</h1>
    {#if data.userData.id === data.user?.id}
      <Button
        variant="ghost"
        size="icon"
        onclick={() => {
          openNewModal = true;
        }}
      >
        <Plus />
      </Button>
    {/if}
  </div>

  <div
    class="grid grid-cols-[repeat(auto-fit,minmax(250px,1fr))] items-center justify-items-center gap-6 p-6"
  >
    {#each data.folders as folder}
      <a
        href="/users/{data.userData.id}/folders/{folder.id}"
        class="w-[240px] cursor-pointer overflow-hidden rounded-xl border shadow-md transition hover:shadow-lg"
      >
        <div class="relative flex h-32 items-center justify-center">
          <div
            class="absolute inset-0 bg-gradient-to-br from-[#F59E0B] to-[#9333EA]"
          ></div>
        </div>
        <div class="flex flex-col justify-between p-4">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-semibold">
              {folder.name}
            </h2>
          </div>
          <p class="mt-1 text-sm text-muted-foreground">
            {folder.itemCount} items
          </p>
        </div>
      </a>
    {/each}
  </div>
</Card.Root>

<!-- <a
        href="/users/{data.userData.id}/folders/{folder.id}"
        class="font-medium hover:underline"
      >
        
      </a> -->

<NewFolderModal bind:open={openNewModal} />
