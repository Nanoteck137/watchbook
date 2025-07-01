<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import {
    MediaRating,
    MediaStatus,
    MediaType,
    parseMediaRating,
    parseMediaStatus,
    parseMediaType,
  } from "$lib/api_types";
  import FormItem from "$lib/components/FormItem.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import { cn } from "$lib/utils";
  import {
    Breadcrumb,
    Button,
    Input,
    Label,
    Select,
  } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";
  import SuperDebug from "sveltekit-superforms";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";

  const apiClient = getApiClient();

  const { data } = $props();

  const mediaTypeLabels: { [K in MediaType]: string } = {
    unknown: "Unknown",
    season: "Season",
    movie: "Movie",
    "anime-season": "Anime Season",
    "anime-movie": "Anime Movie",
  };

  const mediaStatusLabels: { [K in MediaStatus]: string } = {
    unknown: "Unknown",
    airing: "Airing",
    finished: "Finished",
    "not-aired": "Not Aired",
  };

  const mediaRatingLabels: { [K in MediaRating]: string } = {
    unknown: "Unknown",
    "all-ages": "All Ages",
    pg: "PG",
    "pg-13": "PG-13",
    "r-17": "R-17",
    "r-mild-nudity": "R-Mild-Nudity",
    "r-hentai": "R-Hentai",
  };

  const formSchema = z.object({
    type: MediaType,

    tmdbId: z.string(),
    malId: z.string(),
    anilistId: z.string(),

    title: z.string().min(1),
    description: z.string(),

    score: z.number(),
    status: MediaStatus,
    rating: MediaRating,
    airingSeason: z.string(),

    startDate: z.string(),
    endDate: z.string(),

    tags: z.string(),
    studios: z.string(),
  });

  const { form, enhance, errors, submitting } = superForm(
    defaults(zod(formSchema)),
    {
      validators: zod(formSchema),
      SPA: true,
      resetForm: false,
      onUpdate: async ({ form: f }) => {
        console.log(f.data);

        if (f.valid) {
          const res = await apiClient.editMedia(data.media.id, {
            title: f.data.title,
            description: f.data.description,
            type: f.data.type,
            score: f.data.score,
            status: f.data.status,
            rating: f.data.rating,
            airingSeason: f.data.airingSeason,
            startDate: f.data.startDate,
            endDate: f.data.endDate,
            studios: f.data.studios.split(","),
            tags: f.data.tags.split(","),
            tmdbId: f.data.tmdbId,
            malId: f.data.malId,
            anilistId: f.data.anilistId,
          });
          if (!res.success) {
            return handleApiError(res.error);
          }

          // TODO(patrik): Navigate to media item
          toast.success("Successfully updated media item");
          await invalidateAll();
        }
      },
    },
  );

  $effect(() => {
    const media = data.media;

    $form.type = parseMediaType(media.type);

    // $form.tmdbId = media.;
    // $form.malId = data?.malId ?? "";
    // $form.anilistId = "";

    $form.title = media.title;
    $form.description = media.description ?? "";

    $form.score = media.score ?? 0.0;
    $form.status = parseMediaStatus(media.status);
    $form.rating = parseMediaRating(media.rating);
    // $form.airingSeason = media.data?.airingSeason ?? "";

    // $form.startDate = media.startDate ?? "";
    // $form.endDate = data?.endDate ?? "";

    $form.tags = media.tags.join(",") ?? "";
    $form.studios = media.studios.join(",") ?? "";

    // $form.coverUrl = data?.coverImageUrl ?? "";

    // $form.episodeCount = data?.episodeCount ?? 0;
    // $form.generateEpisodes = true;
  });
</script>

<div class="py-2">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/media">Media</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Link href={`/media/${data.media.id}`}>
          {data.media.title}
        </Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page class="line-clamp-1 max-w-96 text-ellipsis">
          Edit
        </Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<Spacer size="sm" />

<form method="POST" class="w-2/3 space-y-6" use:enhance>
  <FormItem>
    <Label for="type" aria-invalid={$errors.type ? "true" : undefined}>
      Type
    </Label>
    <Select.Root
      type="single"
      bind:value={$form.type}
      name="type"
      allowDeselect={false}
    >
      <Select.Trigger id="type">
        {mediaTypeLabels[$form.type]}
      </Select.Trigger>
      <Select.Content>
        {#each MediaType.options as opt}
          <Select.Item value={opt} label={mediaTypeLabels[opt]} />
        {/each}
      </Select.Content>
    </Select.Root>
    {#if $errors.type}
      <span class="invalid">{$errors.type}</span>
    {/if}
  </FormItem>

  <FormItem>
    <Label for="tmdbId" aria-invalid={$errors.tmdbId ? "true" : undefined}>
      TheMovieDB ID
    </Label>
    <Input id="tmdbId" name="tmdbId" type="text" bind:value={$form.tmdbId} />
    {#if $errors.tmdbId}
      <span class="invalid">{$errors.tmdbId}</span>
    {/if}
  </FormItem>

  <FormItem>
    <Label for="malId" aria-invalid={$errors.malId ? "true" : undefined}>
      MyAnimeList ID
    </Label>
    <Input id="malId" name="malId" type="text" bind:value={$form.malId} />
    {#if $errors.malId}
      <span class="invalid">{$errors.malId}</span>
    {/if}
  </FormItem>

  <FormItem>
    <Label
      for="anilistId"
      aria-invalid={$errors.anilistId ? "true" : undefined}
    >
      Anilist ID
    </Label>
    <Input
      id="anilistId"
      name="anilistId"
      type="text"
      bind:value={$form.anilistId}
    />
    {#if $errors.anilistId}
      <span class="invalid">{$errors.anilistId}</span>
    {/if}
  </FormItem>

  <FormItem>
    <Label for="title" aria-invalid={$errors.title ? "true" : undefined}>
      Title*
    </Label>
    <Input id="title" name="title" type="text" bind:value={$form.title} />
    {#if $errors.title}
      <span class="invalid">{$errors.title}</span>
    {/if}
  </FormItem>

  <FormItem>
    <Label
      for="description"
      aria-invalid={$errors.description ? "true" : undefined}
    >
      Description
    </Label>
    <textarea
      id="description"
      class={cn(
        "aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive field-sizing-content shadow-xs flex min-h-16 w-full rounded-md border border-input bg-transparent px-3 py-2 text-base outline-none transition-[color,box-shadow] placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-input/30 md:text-sm",
        "resize-none",
      )}
      value={$form.description}
      rows={8}
    ></textarea>
    {#if $errors.description}
      <span class="invalid">{$errors.description}</span>
    {/if}
  </FormItem>

  <FormItem>
    <Label for="score" aria-invalid={$errors.score ? "true" : undefined}>
      Score
    </Label>
    <Input
      id="score"
      name="score"
      type="number"
      step="0.01"
      bind:value={$form.score}
    />
    {#if $errors.score}
      <span class="invalid">{$errors.score}</span>
    {/if}
  </FormItem>

  <FormItem>
    <Label for="status" aria-invalid={$errors.status ? "true" : undefined}>
      Status
    </Label>
    <Select.Root
      type="single"
      bind:value={$form.status}
      name="status"
      allowDeselect={false}
    >
      <Select.Trigger id="status">
        {mediaStatusLabels[$form.status]}
      </Select.Trigger>
      <Select.Content>
        {#each MediaStatus.options as opt}
          <Select.Item value={opt} label={mediaStatusLabels[opt]} />
        {/each}
      </Select.Content>
    </Select.Root>
    {#if $errors.status}
      <span class="invalid">{$errors.status}</span>
    {/if}
  </FormItem>

  <FormItem>
    <Label for="rating" aria-invalid={$errors.status ? "true" : undefined}>
      Rating
    </Label>
    <Select.Root
      type="single"
      bind:value={$form.rating}
      name="rating"
      allowDeselect={false}
    >
      <Select.Trigger id="rating">
        {mediaRatingLabels[$form.rating]}
      </Select.Trigger>
      <Select.Content>
        {#each MediaRating.options as opt}
          <Select.Item value={opt} label={mediaRatingLabels[opt]} />
        {/each}
      </Select.Content>
    </Select.Root>
    {#if $errors.rating}
      <span class="invalid">{$errors.rating}</span>
    {/if}
  </FormItem>

  <FormItem>
    <Label
      for="airingSeason"
      aria-invalid={$errors.airingSeason ? "true" : undefined}
    >
      Airing Season
    </Label>
    <Input
      id="airingSeason"
      name="airingSeason"
      type="text"
      bind:value={$form.airingSeason}
    />
    {#if $errors.airingSeason}
      <span class="invalid">{$errors.airingSeason}</span>
    {/if}
  </FormItem>

  <FormItem>
    <Label
      for="startDate"
      aria-invalid={$errors.startDate ? "true" : undefined}
    >
      Start Date
    </Label>
    <Input
      id="startDate"
      name="startDate"
      type="text"
      bind:value={$form.startDate}
    />
    {#if $errors.startDate}
      <span class="invalid">{$errors.startDate}</span>
    {/if}
  </FormItem>

  <FormItem>
    <Label for="endDate" aria-invalid={$errors.endDate ? "true" : undefined}>
      End Date
    </Label>
    <Input
      id="endDate"
      name="endDate"
      type="text"
      bind:value={$form.endDate}
    />
    {#if $errors.endDate}
      <span class="invalid">{$errors.endDate}</span>
    {/if}
  </FormItem>

  <FormItem>
    <Label for="tags" aria-invalid={$errors.tags ? "true" : undefined}>
      Tags
    </Label>
    <Input id="tags" name="tags" type="text" bind:value={$form.tags} />
    {#if $errors.tags}
      <span class="invalid">{$errors.tags}</span>
    {/if}
  </FormItem>

  <FormItem>
    <Label for="studios" aria-invalid={$errors.studios ? "true" : undefined}>
      Studios
    </Label>
    <Input
      id="studios"
      name="studios"
      type="text"
      bind:value={$form.studios}
    />
    {#if $errors.studios}
      <span class="invalid">{$errors.studios}</span>
    {/if}
  </FormItem>

  <Button type="submit" disabled={$submitting}>Submit</Button>
</form>

<SuperDebug data={$form} />
