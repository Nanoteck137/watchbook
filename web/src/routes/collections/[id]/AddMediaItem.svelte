<script lang="ts">
  import { getApiClient, handleApiError } from "$lib";
  import FormItem from "$lib/components/FormItem.svelte";
  import Image from "$lib/components/Image.svelte";
  import type { Modal } from "$lib/components/modals";
  import { cn } from "$lib/utils";
  import {
    Button,
    Dialog,
    Input,
    Label,
    ScrollArea,
  } from "@nanoteck137/nano-ui";
  import { Check } from "lucide-svelte";
  import MultiEditMediaItem from "./MultiEditMediaItem.svelte";

  export type Props = {
    open: boolean;
    itemIds: string[];
  };

  let {
    open = $bindable(),
    onResult,

    itemIds,
  }: Props & Modal<MediaItem[]> = $props();

  const apiClient = getApiClient();

  let editModalOpen = $state(false);

  export type MediaItem = {
    id: string;
    name: string;
    searchSlug: string;
    order: number;
  };

  export type SearchResult = {
    id: string;
    imageUrl?: string;
    name: string;
    checked: boolean;
    alreadyAdded: boolean;
  };

  let results = $state<SearchResult[]>([]);
  let checkedItems = $derived(
    results.filter((i) => !i.alreadyAdded).filter((i) => i.checked),
  );

  function debounce<T extends (...args: any[]) => void>(fn: T, delay: number) {
    let timer: number | undefined;
    return (...args: Parameters<T>) => {
      if (timer) clearTimeout(timer);
      timer = window.setTimeout(() => fn(...args), delay);
    };
  }

  async function search(query: string) {
    const media = await apiClient.getMedia({
      query: { filter: `title % \"%${query}%"`, perPage: "20" },
    });
    if (!media.success) {
      return handleApiError(media.error);
    }

    results = media.data.media.map((m) => ({
      id: m.id,
      imageUrl: m.coverUrl ?? undefined,
      name: m.title,
      checked: false,
      alreadyAdded: !!itemIds.find((i) => i === m.id),
    }));
  }

  const debouncedSearch = debounce((query: string) => {
    search(query);
  }, 300);
</script>

<Dialog.Root bind:open>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Add media to the collection</Dialog.Title>
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
                  class="h-full w-full"
                  src={result.imageUrl}
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
                  title={result.name}
                >
                  {result.name}
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
        onclick={() => {
          editModalOpen = true;
        }}
      >
        Add ({checkedItems.length} items)
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>

<MultiEditMediaItem
  bind:open={editModalOpen}
  items={checkedItems}
  onResult={(items) => {
    onResult(items);
    open = false;
  }}
/>
