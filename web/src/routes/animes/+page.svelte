<script lang="ts">
  import Image from "$lib/components/Image.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import { pickTitle } from "$lib/utils.js";
  import { Star } from "lucide-svelte";

  const { data } = $props();

  // NOTE(patrik):
  //  - Large image size: 225x318 (h-80 w-56)
  //  - Medium image size: 160x220 (h-56 w-40)
  //  - Small image size: 50x70 (h-20 w-14)
</script>

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
      <div class="flex items-center justify-center border-l px-4">
        <Star size={18} class="fill-yellow-200 text-yellow-200" />
        <Spacer horizontal size="xs" />
        <p class="text-xs">8.82</p>
      </div>
    </div>
  {/each}
</div>
