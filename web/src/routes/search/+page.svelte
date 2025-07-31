<script lang="ts">
  import { goto, invalidateAll, onNavigate } from "$app/navigation";
  import Image from "$lib/components/Image.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import { Button, Input, Label } from "@nanoteck137/nano-ui";
  import { Star } from "lucide-svelte";

  const { data } = $props();

  async function search(query: string) {
    await goto(`?query=${query}`, {
      invalidateAll: true,
      keepFocus: true,
      replaceState: true,
    });
  }

  let value = "";

  let timer: NodeJS.Timeout;
  function onInput(e: Event) {
    const target = e.target as HTMLInputElement;
    const current = target.value;
    value = current;

    clearTimeout(timer);
    timer = setTimeout(async () => {
      search(current);
    }, 500);
  }

  // NOTE(patrik): Fix for clicking the search button
  onNavigate((e) => {
    if (e.type === "link" && e.from?.url.pathname === "/search") {
      invalidateAll();
    }
  });

  function formatError(err: { type: string; code: number; message: string }) {
    // TODO(patrik): Better error
    return err.message;
  }
</script>

<form
  action=""
  method="get"
  onsubmit={(e) => {
    e.preventDefault();
    clearTimeout(timer);
    search(value);
  }}
>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-2">
      <Label for="query">Search Query</Label>
      <Input
        id="query"
        name="query"
        autocomplete="off"
        value={data.query}
        oninput={onInput}
      />
    </div>
    <Button type="submit">Search</Button>
  </div>
</form>

<div class="h-4"></div>

{#if data.mediaError}
  <p class="text-red-400">
    Media Query Error: {data.mediaError}
  </p>
{/if}

{#if data.media.length > 0}
  <div class="flex items-center justify-between">
    <p class="text-bold">Media</p>
    <p class="text-xs">{data.media.length} Entries</p>
  </div>

  {#each data.media as media}
    <div class="flex justify-between border-b py-2">
      <div class="flex">
        <Image class="h-20 w-14" src={media.coverUrl} alt="cover" />
        <div class="px-4 py-1">
          <a
            class="line-clamp-2 text-ellipsis text-sm font-semibold hover:cursor-pointer hover:underline"
            href="/media/{media.id}"
            title={media.title}
          >
            {media.title}
          </a>
        </div>
      </div>
      <div class="flex min-w-24 max-w-24 items-center justify-center border-l">
        <Star size={18} class="fill-foreground" />
        <Spacer horizontal size="xs" />
        <p class="font-mono text-xs">{media.score?.toFixed(2) ?? "N/A"}</p>
      </div>
    </div>
  {/each}
{/if}
