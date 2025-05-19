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
  import ImportMalListDialog from "./ImportMalListDialog.svelte";
  import ImportAnime from "./ImportAnime.svelte";

  const { data } = $props();
  const apiClient = getApiClient();

  // NOTE(patrik):
  //  - Large image size: 225x318 (h-80 w-56)
  //  - Medium image size: 160x220 (h-56 w-40)
  //  - Small image size: 50x70 (h-20 w-14)
</script>

<ImportMalListDialog />
<ImportAnime />

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
