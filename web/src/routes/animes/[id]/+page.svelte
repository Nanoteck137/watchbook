<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import { Anime } from "$lib/api/types.js";
  import Image from "$lib/components/Image.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import { cn, pickTitle } from "$lib/utils.js";
  import {
    Breadcrumb,
    Button,
    buttonVariants,
    Checkbox,
    Dialog,
    DropdownMenu,
    Input,
    Label,
    Separator,
  } from "@nanoteck137/nano-ui";
  import { ChevronDown, Eye, Star } from "lucide-svelte";
  import toast from "svelte-5-french-toast";
  import { z } from "zod";

  const { data } = $props();
  const apiClient = getApiClient();

  let showMore = $state(false);
  let episodeOpen = $state(false);

  function formatAnimeType(ty: string) {
    switch (ty) {
      case "tv":
        return "TV";
    }

    return ty;
  }

  async function updateScore(score: number | null) {
    if (score === null) {
      const res = await apiClient.setAnimeUserData(data.anime.id, {
        score: 0,
      });
      if (!res.success) {
        return handleApiError(res.error);
      }
    } else {
      const res = await apiClient.setAnimeUserData(data.anime.id, { score });
      if (!res.success) {
        return handleApiError(res.error);
      }
    }

    await invalidateAll();
  }

  type List =
    | "watching"
    | "completed"
    | "on-hold"
    | "dropped"
    | "plan-to-watch";

  async function updateList(list: List | null) {
    if (list === null) {
      const res = await apiClient.setAnimeUserData(data.anime.id, {
        list: "",
      });
      if (!res.success) {
        return handleApiError(res.error);
      }
    } else {
      const res = await apiClient.setAnimeUserData(data.anime.id, {
        list,
      });
      if (!res.success) {
        return handleApiError(res.error);
      }
    }

    await invalidateAll();
  }

  function formatScore() {
    const score = data.anime.user?.score;

    if (score === null) return "-";

    return score?.toString();
  }

  function formatList() {
    const category = data.anime.user?.list;

    if (!category) return "Not Added";

    switch (category) {
      case "watching":
        return "Watching";
      case "completed":
        return "Completed";
      case "on-hold":
        return "On-Hold";
      case "dropped":
        return "Dropped";
      case "plan-to-watch":
        return "Plan to Watch";
    }

    return category;
  }
</script>

<div class="py-2">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/animes">Animes</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page class="line-clamp-1 max-w-96 text-ellipsis"
          >{pickTitle(data.anime)}</Breadcrumb.Page
        >
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<Spacer size="sm" />

<div class="flex flex-col items-center gap-2">
  <Image class="min-h-80 w-56" src={data.anime.coverUrl} alt="cover" />
  <div class="flex flex-col gap-1">
    <p class="text-center text-base">{data.anime.title}</p>
    {#if data.anime.titleEnglish}
      <Separator />
      <p class="text-center text-sm text-zinc-300">
        {data.anime.titleEnglish}
      </p>
    {/if}
  </div>
</div>

<Spacer size="lg" />

<div class="flex flex-col justify-around gap-2">
  <DropdownMenu.Root>
    <DropdownMenu.Trigger
      class="relative flex items-center justify-center rounded bg-primary py-1 text-primary-foreground"
    >
      <Star size={18} class="fill-primary-foreground" />
      <Spacer horizontal size="sm" />
      <p class="text-base">{formatList()}</p>
      <ChevronDown class="absolute right-4" size={20} />
    </DropdownMenu.Trigger>
    <DropdownMenu.Content class="w-40">
      <DropdownMenu.Group>
        <DropdownMenu.GroupHeading>Select Category</DropdownMenu.GroupHeading>
        <DropdownMenu.Separator />
        <DropdownMenu.Item onclick={() => updateList("watching")}>
          Watching
        </DropdownMenu.Item>
        <DropdownMenu.Item onclick={() => updateList("completed")}>
          Completed
        </DropdownMenu.Item>
        <DropdownMenu.Item onclick={() => updateList("on-hold")}>
          On-Hold
        </DropdownMenu.Item>
        <DropdownMenu.Item onclick={() => updateList("dropped")}>
          Dropped
        </DropdownMenu.Item>
        <DropdownMenu.Item onclick={() => updateList("plan-to-watch")}>
          Plan to Watch
        </DropdownMenu.Item>
        {#if !!data.anime.user?.list}
          <DropdownMenu.Separator />
          <DropdownMenu.Item onclick={() => updateList(null)}>
            Remove
          </DropdownMenu.Item>
        {/if}
      </DropdownMenu.Group>
    </DropdownMenu.Content>
  </DropdownMenu.Root>

  <DropdownMenu.Root>
    <DropdownMenu.Trigger
      class="relative flex items-center justify-center rounded bg-primary py-1 text-primary-foreground"
    >
      <Star size={18} class="fill-primary-foreground" />
      <Spacer horizontal size="sm" />
      <p class="text-base">{formatScore()}</p>
      <ChevronDown class="absolute right-4" size={20} />
    </DropdownMenu.Trigger>
    <DropdownMenu.Content class="w-40">
      <DropdownMenu.Group>
        <DropdownMenu.GroupHeading>Select Rating</DropdownMenu.GroupHeading>
        <DropdownMenu.Separator />
        <DropdownMenu.Item onclick={() => updateScore(10)}>
          (10) Masterpiece
        </DropdownMenu.Item>
        <DropdownMenu.Item onclick={() => updateScore(9)}>
          (9) Great
        </DropdownMenu.Item>
        <DropdownMenu.Item onclick={() => updateScore(8)}>
          (8) Very Good
        </DropdownMenu.Item>
        <DropdownMenu.Item onclick={() => updateScore(7)}>
          (7) Good
        </DropdownMenu.Item>
        <DropdownMenu.Item onclick={() => updateScore(6)}>
          (6) Fine
        </DropdownMenu.Item>
        <DropdownMenu.Item onclick={() => updateScore(5)}>
          (5) Average
        </DropdownMenu.Item>
        <DropdownMenu.Item onclick={() => updateScore(4)}>
          (4) Bad
        </DropdownMenu.Item>
        <DropdownMenu.Item onclick={() => updateScore(3)}>
          (3) Very Bad
        </DropdownMenu.Item>
        <DropdownMenu.Item onclick={() => updateScore(2)}>
          (2) Horrible
        </DropdownMenu.Item>
        <DropdownMenu.Item onclick={() => updateScore(1)}>
          (1) Appalling
        </DropdownMenu.Item>
        <DropdownMenu.Item onclick={() => updateScore(null)}>
          No Score
        </DropdownMenu.Item>
      </DropdownMenu.Group>
    </DropdownMenu.Content>
  </DropdownMenu.Root>

  <!-- <Separator orientation="vertical" /> -->

  <Dialog.Root bind:open={episodeOpen}>
    <Dialog.Trigger class={cn(buttonVariants(), "w-full")}>
      <div class="flex items-center justify-center rounded">
        <Eye size={18} />
        <Spacer horizontal size="sm" />
        <p class="text-base">
          {data.anime.user?.episode ?? "??"} / {data.anime.episodeCount ??
            "??"}
        </p>
      </div>
    </Dialog.Trigger>
    <Dialog.Content class="sm:max-w-[425px]">
      <form
        class="w-full"
        onsubmit={async (e) => {
          e.preventDefault();

          const formData = new FormData(e.target as HTMLFormElement);
          console.log(formData);

          const episode = formData.get("episode")?.toString() ?? "0";
          const isRewatching = formData.get("isRewatching")?.toString() ?? "";

          const res = await apiClient.setAnimeUserData(data.anime.id, {
            episode: parseInt(episode === "" ? "0" : episode),
            isRewatching: isRewatching === "on",
          });

          if (!res.success) {
            return handleApiError(res.error);
          }

          toast.success("Successfully update episode");
          invalidateAll();
          episodeOpen = false;
        }}
      >
        <Dialog.Header>
          <Dialog.Title>Edit episode</Dialog.Title>
          <!-- <Dialog.Description>
          Make changes to your profile here. Click save when you're done.
        </Dialog.Description> -->
        </Dialog.Header>
        <div class="grid gap-4 py-4">
          <div class="flex flex-col gap-2">
            <Label for="episode">Current episode</Label>
            <Input
              name="episode"
              id="episode"
              type="number"
              value={data.anime.user?.episode ?? 0}
            />
          </div>
          <div class="flex items-center gap-2">
            <Checkbox
              name="isRewatching"
              id="isRewatching"
              checked={data.anime.user?.isRewatching ?? false}
            />
            <Label for="isRewatching">Is Rewatching</Label>
          </div>
        </div>
        <Dialog.Footer>
          <Button type="submit">Save changes</Button>
        </Dialog.Footer>
      </form>
    </Dialog.Content>
  </Dialog.Root>
  <!-- <Button class=""></Button> -->
</div>

<Spacer size="lg" />

<!-- {#if data.anime.description}
  <div class="flex flex-col gap-1">
    <p
      class={`text-ellipsis whitespace-pre-line text-sm ${!showMore ? "line-clamp-2" : ""}`}
    >
      {data.anime.description}
    </p>
    <Button
      class="w-fit"
      size="sm"
      variant="outline"
      onclick={() => {
        showMore = !showMore;
      }}
    >
      Show More
    </Button>
  </div>
{/if} -->

<!-- <Spacer size="lg" /> -->

<div
  class="flex items-center justify-center rounded bg-primary py-1 text-primary-foreground"
>
  <Star size={18} />
  <Spacer horizontal size="xs" />
  <span>Score:</span>
  <Spacer horizontal size="xs" />
  <p class="text-base">{data.anime.score?.toFixed(2) ?? "N/A"}</p>
</div>

<Spacer size="xs" />

<div
  class="flex flex-col gap-1 rounded bg-primary p-2 text-primary-foreground"
>
  <p>Type: {formatAnimeType(data.anime.type)}</p>
  <p>Episodes: {data.anime.episodeCount}</p>
  <p>Status: {data.anime.status}</p>
  <p>Start Date: {data.anime.startDate ?? "Unknown"}</p>
  <p>End Date: {data.anime.endDate ?? "Unknown"}</p>
  <p>
    Studios:
    {#each data.anime.studios as studio, i}
      {#if i != 0}
        <span>, </span>
      {/if}
      <a class="text-blue-500 hover:underline" href="/studios/{studio.slug}">
        {studio.name}
      </a>
    {/each}
  </p>
  <p>
    Producers:
    {#each data.anime.producers as producer, i}
      {#if i != 0}
        <span>, </span>
      {/if}
      <a
        class="text-blue-500 hover:underline"
        href="/producers/{producer.slug}"
      >
        {producer.name}
      </a>
    {/each}
  </p>
  <p>
    Genres:
    {#each data.anime.genres as genre, i}
      {#if i != 0}
        <span>, </span>
      {/if}
      <a class="text-blue-500 hover:underline" href="/tags/{genre.slug}">
        {genre.name}
      </a>
    {/each}
  </p>
  <p>
    Themes:
    {#each data.anime.themes as theme, i}
      {#if i != 0}
        <span>, </span>
      {/if}
      <a class="text-blue-500 hover:underline" href="/tags/{theme.slug}">
        {theme.name}
      </a>
    {/each}
  </p>
  <p>
    Demographics:

    {#each data.anime.demographics as demographic, i}
      {#if i != 0}
        <span>, </span>
      {/if}
      <a class="text-blue-500 hover:underline" href="/tags/{demographic.slug}">
        {demographic.name}
      </a>
    {/each}
  </p>
  <p>Rating: {data.anime.rating}</p>
</div>
