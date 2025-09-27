<script lang="ts">
  import { Button } from "@nanoteck137/nano-ui";
  import NewFolderModal from "./NewFolderModal.svelte";

  const { data } = $props();

  let openNewModal = $state(false);
</script>

<Button
  onclick={() => {
    openNewModal = true;
  }}>New</Button
>

{#each data.folders as folder}
  <div class="flex items-center justify-between border-b py-4">
    <div>
      <a
        href="/users/{data.userData.id}/folders/{folder.id}"
        class="font-medium hover:underline"
      >
        {folder.name}
      </a>
      <!-- <p class="text-xs text-muted-foreground">Air Date: {part.releaseDate}</p> -->
    </div>

    <div class="flex items-center gap-2">
      <Button variant="outline" size="sm" onclick={() => {}}>Edit</Button>
      <Button variant="destructive" size="sm" onclick={() => {}}>
        Delete
      </Button>
    </div>
  </div>
{/each}

<NewFolderModal bind:open={openNewModal} />
