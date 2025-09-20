<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import { mediaRatings, mediaStatus } from "$lib/api-types.js";
  import Spacer from "$lib/components/Spacer.svelte";
  import type { UserList } from "$lib/types.js";
  import { Button, Card } from "@nanoteck137/nano-ui";

  const { data } = $props();
  const apiClient = getApiClient();

  let episodeOpen = $state(false);
  let descriptionShowMore = $state(false);

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

  function formatDate(dateString?: string) {
    if (!dateString) return "N/A"; // handle missing dates

    const date = new Date(dateString);

    // const locale = "sv-SE";
    // const locale = "en-US";
    const locale = "en-GB";
    return new Intl.DateTimeFormat(locale, {
      year: "numeric",
      month: "long", // "short" → Apr, "numeric" → 4
      day: "numeric",
    }).format(date);
  }
</script>

<Card.Root>
  <Card.Header>
    <!-- <Card.Title>Card Title</Card.Title>
    <Card.Description>Card Description</Card.Description> -->

    <Card.Title class="text-xl">Description</Card.Title>
    <div class="flex flex-col gap-1">
      <p
        class={`text-ellipsis whitespace-pre-line text-sm ${!descriptionShowMore ? "line-clamp-4" : ""}`}
      >
        {data.media.description}
      </p>

      <Button
        class="w-fit"
        size="sm"
        variant="outline"
        onclick={() => {
          descriptionShowMore = !descriptionShowMore;
        }}
      >
        Show More
      </Button>
    </div>
  </Card.Header>

  <Card.Content>
    <Card.Title class="text-xl">Information</Card.Title>
    <Spacer />
    <!-- <Card.Description>Card Description</Card.Description> -->
    <dl
      class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4"
    >
      <div>
        <dt class="font-medium">Score</dt>
        <dd class="text-muted-foreground">
          {data.media.score ?? "N/A"}
        </dd>
      </div>

      <div>
        <dt class="font-medium">Parts</dt>
        <dd class="text-muted-foreground">{data.media.partCount}</dd>
      </div>

      <div>
        <dt class="font-medium">Status</dt>
        <dd class="text-muted-foreground">
          {mediaStatus.find((i) => i.value === data.media.status)?.label}
        </dd>
      </div>

      <div>
        <dt class="font-medium">Airing Season</dt>
        <dd class="text-muted-foreground">
          {data.media.airingSeason ?? "N/A"}
        </dd>
      </div>

      <div>
        <dt class="font-medium">Start Date</dt>
        <dd class="text-muted-foreground">
          {formatDate(data.media.startDate ?? undefined)}
        </dd>
      </div>

      <div>
        <dt class="font-medium">End Date</dt>
        <dd class="text-muted-foreground">
          {formatDate(data.media.endDate ?? undefined)}
        </dd>
      </div>

      <div>
        <dt class="font-medium">Rating</dt>
        <dd class="text-muted-foreground">
          {mediaRatings.find((i) => i.value === data.media.rating)?.label}
        </dd>
      </div>

      <div class="sm:col-span-2 md:col-span-3 lg:col-span-4">
        <dt class="font-medium">Creators</dt>

        <dd class="mt-1 flex flex-wrap gap-2">
          {#each data.media.creators as creator}
            <span
              class="rounded-md bg-gray-100 px-2 py-1 text-xs text-gray-700"
            >
              {creator}
            </span>
          {/each}
        </dd>
      </div>

      <div class="sm:col-span-2 md:col-span-3 lg:col-span-4">
        <dt class="font-medium">Tags</dt>

        <dd class="mt-1 flex flex-wrap gap-2">
          {#each data.media.tags as tag}
            <a
              class="rounded-md bg-gray-100 px-2 py-1 text-xs text-gray-700"
              href="/media?filterTags={tag}"
            >
              {tag}
            </a>
          {/each}
        </dd>
      </div>
    </dl>
  </Card.Content>
</Card.Root>
