<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import type { Media } from "$lib/api/types";
  import FormItem from "$lib/components/FormItem.svelte";
  import Image from "$lib/components/Image.svelte";
  import { cn, debounce } from "$lib/utils";
  import {
    Button,
    Dialog,
    Input,
    Label,
    ScrollArea,
  } from "@nanoteck137/nano-ui";
  import { Check } from "lucide-svelte";
  import toast from "svelte-5-french-toast";

  export type Props = {
    open: boolean;
    folderId: string;
    itemIds: string[];
  };

  let { open = $bindable(), folderId, itemIds }: Props = $props();

  const apiClient = getApiClient();

  export type Item = {
    media: Media;
    checked: boolean;
    alreadyAdded: boolean;
  };

  let results = $state<Item[]>([]);
  let checkedItems = $derived(
    results.filter((i) => !i.alreadyAdded).filter((i) => i.checked),
  );

  async function search(query: string) {
    const media = await apiClient.getMedia({
      query: { filter: `title % \"%${query}%"`, perPage: "20" },
    });
    if (!media.success) {
      return handleApiError(media.error);
    }

    results = media.data.media.map((m) => ({
      media: m,
      alreadyAdded: !!itemIds.find((i) => i === m.id),
      checked: false,
    }));
  }

  const debouncedSearch = debounce((query: string) => {
    search(query);
  }, 300);

  $effect(() => {
    if (open) {
      results = [];
    }
  });
</script>

<Dialog.Root bind:open>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Add media to this folder</Dialog.Title>
    </Dialog.Header>

    <FormItem>
      <Label for="search">Search</Label>
      <Input
        id="search"
        name="search"
        type="text"
        oninput={(e) => {
          const val = (e.target as HTMLInputElement).value;
          debouncedSearch(val);
        }}
      />
      <!-- <Errors errors={$errors.name} /> -->
    </FormItem>

    <ScrollArea class="h-60">
      <Label>Results</Label>

      <div class="flex flex-col gap-2">
        {#each results as result, i}
          <button
            class="group flex justify-between border-b py-2"
            disabled={result.alreadyAdded}
            onclick={() => {
              results[i].checked = !results[i].checked;
            }}
          >
            <div class="flex">
              <div class="relative h-20 w-14">
                <Image
                  class="min-h-20 min-w-14"
                  src={result.media.coverUrl}
                  alt="cover"
                />
                <div
                  class={cn(
                    "absolute inset-0 flex items-center justify-center ",
                    result.alreadyAdded
                      ? "bg-black/40"
                      : "group-hover:bg-black/40",
                  )}
                >
                  {#if result.checked || result.alreadyAdded}
                    <div
                      class="flex h-8 w-8 items-center justify-center rounded-full bg-black/80"
                    >
                      <Check />
                    </div>
                  {/if}
                </div>
              </div>
              <div class="flex flex-col gap-2 px-4 py-1">
                <p
                  class="line-clamp-2 text-ellipsis text-sm font-semibold"
                  title={result.media.title}
                >
                  {result.media.title}
                </p>
                <!-- <p class="text-xs">{media.user?.list}</p> -->
              </div>
            </div>
            <!-- <div
              class="flex min-w-24 max-w-24 items-center justify-center border-l"
            >
              <Star size={18} class="fill-foreground" />
              <Spacer horizontal size="xs" />
              <p class="font-mono text-xs">
                {media.user?.score ?? "0"}
              </p>
            </div> -->
          </button>
        {/each}
      </div>
    </ScrollArea>

    <Dialog.Footer class="gap-2 sm:gap-0">
      <Button
        variant="outline"
        onclick={() => {
          open = false;
        }}
      >
        Close
      </Button>

      <Button
        disabled={checkedItems.length <= 0}
        onclick={async () => {
          console.log(checkedItems);
          let error = false;
          for (const item of checkedItems) {
            const res = await apiClient.addFolderItem(folderId, item.media.id);
            if (!res.success) {
              handleApiError(res.error);
              error = true;
              continue;
            }
          }

          if (error) {
            toast.error("Some items failed to be added");
          } else {
            toast.success("Successfully added items to folder");
          }

          invalidateAll();
          open = false;
        }}
      >
        Add ({checkedItems.length} items)
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
