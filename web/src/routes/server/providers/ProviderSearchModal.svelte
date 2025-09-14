<script lang="ts">
  import { getApiClient, handleApiError } from "$lib";
  import { ProviderSearchResult } from "$lib/api/types";
  import FormItem from "$lib/components/FormItem.svelte";
  import Image from "$lib/components/Image.svelte";
  import type { Modal } from "$lib/components/modals";
  import { cn, debounce } from "$lib/utils";
  import {
    Button,
    Dialog,
    Input,
    Label,
    ScrollArea,
  } from "@nanoteck137/nano-ui";
  import { Search } from "lucide-svelte";

  export type Props = {
    providerName: string;
  };

  const {
    providerName,

    class: className,
    children,
    onResult,
  }: Props & Modal<void> = $props();

  const apiClient = getApiClient();

  let open = $state(false);

  let results = $state<ProviderSearchResult[]>([]);

  async function search(query: string) {
    const res = await apiClient.providerSearchMedia(providerName, {
      query: { query },
    });
    if (!res.success) {
      return handleApiError(res.error);
    }

    console.log(res.data);

    results = res.data.searchResults;

    // results = media.data.media.map((m) => ({
    //   id: m.id,
    //   imageUrl: m.coverUrl ?? undefined,
    //   name: m.title,
    //   checked: false,
    //   alreadyAdded: !!itemIds.find((i) => i === m.id),
    // }));
  }

  const debouncedSearch = debounce((query: string) => {
    search(query);
  }, 300);
</script>

<Dialog.Root bind:open>
  <Dialog.Trigger class={className}>
    {@render children?.()}
  </Dialog.Trigger>

  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Search media using '{providerName}' provider</Dialog.Title>
    </Dialog.Header>

    <form
      onsubmit={(e) => {
        e.preventDefault();
        const data = new FormData(e.target as HTMLFormElement);
        const searchQuery = data.get("search");
        if (searchQuery) {
          search(searchQuery.toString());
        }
      }}
    >
      <FormItem>
        <Label for="search">Search</Label>
        <div class="flex items-center gap-2">
          <Input id="search" name="search" type="text" />

          <Button variant="ghost" size="icon" type="submit">
            <Search />
          </Button>
        </div>
        <!-- <Errors errors={$errors.name} /> -->
      </FormItem>
    </form>

    <ScrollArea class="h-60">
      <Label>Results</Label>

      <div class="flex flex-col gap-2">
        {#each results as result, i}
          <button class="group flex justify-between border-b py-2">
            <div class="flex">
              <div class="relative h-20 w-14">
                <Image
                  class="h-full w-full"
                  src={result.imageUrl}
                  alt="cover"
                />
                <!-- <div class="absolute inset-0 flex items-center justify-center">
                  {#if result.checked || result.alreadyAdded}
                    <div
                      class="flex h-8 w-8 items-center justify-center rounded-full bg-black/80"
                    >
                      <Check />
                    </div>
                  {/if}
                </div> -->
              </div>
              <div class="flex flex-col gap-2 px-4 py-1">
                <p
                  class="line-clamp-2 text-ellipsis text-sm font-semibold"
                  title={result.title}
                >
                  {result.title}
                </p>
                <p class="text-start text-xs">ID: {result.providerId}</p>
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
        onclick={() => {
          open = false;
        }}
      >
        Add
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
