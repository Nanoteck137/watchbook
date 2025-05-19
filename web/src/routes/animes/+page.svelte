<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import Image from "$lib/components/Image.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import { cn, pickTitle } from "$lib/utils.js";
  import {
    Button,
    buttonVariants,
    Checkbox,
    Dialog,
    Input,
    Label,
  } from "@nanoteck137/nano-ui";
  import { Eye, Star } from "lucide-svelte";
  import toast from "svelte-5-french-toast";

  const { data } = $props();
  const apiClient = getApiClient();

  let importMalListDialogOpen = $state(false);

  // NOTE(patrik):
  //  - Large image size: 225x318 (h-80 w-56)
  //  - Medium image size: 160x220 (h-56 w-40)
  //  - Small image size: 50x70 (h-20 w-14)
</script>

<Dialog.Root bind:open={importMalListDialogOpen}>
  <Dialog.Trigger class={cn(buttonVariants())}>Import MAL List</Dialog.Trigger>
  <Dialog.Content class="sm:max-w-[425px]">
    <form
      class="w-full"
      onsubmit={async (e) => {
        e.preventDefault();

        const formData = new FormData(e.target as HTMLFormElement);
        console.log(formData);

        const username = formData.get("username")?.toString() ?? "0";
        const overrideExistingEntries =
          formData.get("overrideExistingEntries")?.toString() ?? "";

        const res = await apiClient.importMalList({
          username: username,
          overrideExistingEntries: overrideExistingEntries === "on",
        });
        if (!res.success) {
          return handleApiError(res.error);
        }

        toast.success("Successfully update episode");
        invalidateAll();
        importMalListDialogOpen = false;
      }}
    >
      <Dialog.Header>
        <Dialog.Title>Import MyAnimeList user list</Dialog.Title>
        <!-- <Dialog.Description>
          Make changes to your profile here. Click save when you're done.
        </Dialog.Description> -->
      </Dialog.Header>
      <div class="grid gap-4 py-4">
        <div class="flex flex-col gap-2">
          <Label for="username">MyAnimeList username</Label>
          <Input name="username" id="username" />
        </div>
        <div class="flex items-center gap-2">
          <Checkbox
            name="overrideExistingEntries"
            id="overrideExistingEntries"
            checked={false}
          />
          <Label for="overrideExistingEntries">Override Existing Entries</Label
          >
        </div>
      </div>
      <Dialog.Footer>
        <Button
          variant="ghost"
          onclick={() => {
            importMalListDialogOpen = false;
          }}
        >
          Close
        </Button>
        <Button type="submit">Save changes</Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>

<div class="flex flex-col gap-4">
  {#each data.animes as anime}
    {@const title = pickTitle(anime)}

    <div class="flex justify-between border-b py-2">
      <div class="flex">
        <Image class="h-20 w-14" src={anime.coverUrl} alt="cover" />
        <div class="px-4 py-1">
          <a
            class="line-clamp-2 text-ellipsis text-sm font-semibold hover:cursor-pointer hover:underline"
            href="/animes/{anime.id}"
            {title}
          >
            {title}
          </a>
        </div>
      </div>
      <div class="flex min-w-24 max-w-24 items-center justify-center border-l">
        <Star size={18} class="fill-foreground" />
        <Spacer horizontal size="xs" />
        <p class="font-mono text-xs">{anime.score?.toFixed(2) ?? "N/A"}</p>
      </div>
    </div>
  {/each}
</div>
