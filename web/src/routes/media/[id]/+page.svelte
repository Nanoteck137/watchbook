<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import Spacer from "$lib/components/Spacer.svelte";
  import type { UserList } from "$lib/types.js";
  import { cn } from "$lib/utils.js";
  import {
    Breadcrumb,
    Button,
    buttonVariants,
    Card,
    Checkbox,
    Dialog,
    DropdownMenu,
    Input,
    Label,
  } from "@nanoteck137/nano-ui";
  import { ChevronDown, Eye, Star, Trash } from "lucide-svelte";
  import toast from "svelte-5-french-toast";
  import BannerHeader from "$lib/components/BannerHeader.svelte";
  import MediaDropdown from "./MediaDropdown.svelte";

  const { data } = $props();
  const apiClient = getApiClient();

  let episodeOpen = $state(false);

  async function updateScore(score: number | null) {
    if (score === null) {
      const res = await apiClient.setMediaUserData(data.media.id, {
        score: 0,
      });
      if (!res.success) {
        return handleApiError(res.error);
      }
    } else {
      const res = await apiClient.setMediaUserData(data.media.id, { score });
      if (!res.success) {
        return handleApiError(res.error);
      }
    }

    await invalidateAll();
  }

  async function updateList(list: UserList | null) {
    if (list === null) {
      const res = await apiClient.setMediaUserData(data.media.id, {
        list: "",
      });
      if (!res.success) {
        return handleApiError(res.error);
      }
    } else {
      const res = await apiClient.setMediaUserData(data.media.id, {
        list,
      });
      if (!res.success) {
        return handleApiError(res.error);
      }
    }

    await invalidateAll();
  }

  function formatScore() {
    const score = data.media.user?.score;

    if (score === null) return "-";

    return score?.toString();
  }

  function formatList() {
    const category = data.media.user?.list;

    if (!category) return "Not Added";

    switch (category) {
      case "in-progress":
        return "In-Progress";
      case "completed":
        return "Completed";
      case "on-hold":
        return "On-Hold";
      case "dropped":
        return "Dropped";
      case "backlog":
        return "Backlog";
    }

    return category;
  }
</script>

<div class="h-[400px]"></div>

<div class="flex flex-col justify-around gap-2">
  <div class="flex gap-2">
    {#if !!data.media.user?.hasData}
      <Button
        onclick={async () => {
          const res = await apiClient.deleteMediaUserData(data.media.id);
          if (!res.success) {
            return handleApiError(res.error);
          }

          invalidateAll();
        }}
      >
        <Trash />
      </Button>
    {/if}

    <DropdownMenu.Root>
      <DropdownMenu.Trigger
        class="relative flex flex-1 items-center justify-center rounded bg-primary py-1 text-primary-foreground"
      >
        <Star size={18} class="fill-primary-foreground" />
        <Spacer horizontal size="sm" />
        <p class="text-base">{formatList()}</p>
        <ChevronDown class="absolute right-4" size={20} />
      </DropdownMenu.Trigger>
      <DropdownMenu.Content class="w-40">
        <DropdownMenu.Group>
          <DropdownMenu.GroupHeading>Select Category</DropdownMenu.GroupHeading
          >
          <DropdownMenu.Separator />
          <DropdownMenu.Item onclick={() => updateList("in-progress")}>
            In-Progress
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
          <DropdownMenu.Item onclick={() => updateList("backlog")}>
            Backlog
          </DropdownMenu.Item>
          {#if !!data.media.user?.list}
            <DropdownMenu.Separator />
            <DropdownMenu.Item onclick={() => updateList(null)}>
              Remove
            </DropdownMenu.Item>
          {/if}
        </DropdownMenu.Group>
      </DropdownMenu.Content>
    </DropdownMenu.Root>
  </div>

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
          {data.media.user?.currentPart ?? "??"} / {data.media.partCount ??
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

          // TODO(patrik): Rename
          const episode = formData.get("episode")?.toString() ?? "0";
          const rewatchCount = formData.get("rewatchCount")?.toString() ?? "0";
          const isRewatching = formData.get("isRewatching")?.toString() ?? "";

          const res = await apiClient.setMediaUserData(data.media.id, {
            currentPart: parseInt(episode === "" ? "0" : episode),
            revisitCount: parseInt(rewatchCount === "" ? "0" : rewatchCount),
            isRevisiting: isRewatching === "on",
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
              value={data.media.user?.currentPart ?? 0}
            />
          </div>
          <div class="flex flex-col gap-2">
            <Label for="rewatchCount">Rewatch Count</Label>
            <Input
              name="rewatchCount"
              id="rewatchCount"
              type="number"
              value={data.media.user?.revisitCount ?? 0}
            />
          </div>
          <div class="flex items-center gap-2">
            <Checkbox
              name="isRewatching"
              id="isRewatching"
              checked={data.media.user?.isRevisiting ?? false}
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

<div
  class="flex items-center justify-center rounded bg-primary py-1 text-primary-foreground"
>
  <Star size={18} />
  <Spacer horizontal size="xs" />
  <span>Score:</span>
  <Spacer horizontal size="xs" />
  <p class="text-base">{data.media.score?.toFixed(2) ?? "N/A"}</p>
</div>

<Spacer size="xs" />

<div
  class="flex flex-col gap-1 rounded bg-primary p-2 text-primary-foreground"
>
  <p>Type: {data.media.mediaType}</p>
  <p>Parts: {data.media.partCount}</p>
  <p>Status: {data.media.status}</p>
  {#if data.media.airingSeason}
    <p>
      Airing Season:
      <a
        class="text-blue-500 hover:underline"
        href="/airing/{data.media.airingSeason}"
      >
        {data.media.airingSeason}
      </a>
    </p>
  {/if}
  {#if data.media.startDate}
    <p>Start Date: {data.media.startDate}</p>
  {/if}
  {#if data.media.endDate}
    <p>End Date: {data.media.endDate}</p>
  {/if}
  <p>
    Creators:
    {#each data.media.creators as creator, i}
      {#if i != 0}
        <span>, </span>
      {/if}
      <a class="text-blue-500 hover:underline" href="/creators/{creator}">
        {creator}
      </a>
    {/each}
  </p>
  <p>
    Tags:
    {#each data.media.tags as tag, i}
      {#if i != 0}
        <span>, </span>
      {/if}
      <a class="text-blue-500 hover:underline" href="/tags/{tag}">
        {tag}
      </a>
    {/each}
  </p>
  <p>Rating: {data.media.rating}</p>
</div>
