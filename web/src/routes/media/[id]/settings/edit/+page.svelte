<script lang="ts">
  import { getApiClient, handleApiError } from "$lib";
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Input, Label, Select } from "@nanoteck137/nano-ui";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";
  import { isMatch } from "date-fns";
  import Spinner from "$lib/components/Spinner.svelte";
  import toast from "svelte-5-french-toast";
  import { invalidateAll } from "$app/navigation";
  import {
    MediaRatingEnum,
    mediaRatings,
    mediaStatus,
    MediaStatusEnum,
    MediaTypeEnum,
    mediaTypes,
    type MediaRating,
    type MediaStatus,
    type MediaType,
  } from "$lib/api-types.js";

  const Schema = z.object({
    type: MediaTypeEnum.default("unknown"),

    title: z.string().min(1, "Title cannot be empty"),
    description: z.string().trim(),

    score: z.number().min(0).max(10),
    status: MediaStatusEnum.default("unknown"),
    rating: MediaRatingEnum.default("unknown"),
    airingSeason: z.string(),

    startDate: z
      .string()
      .trim()
      .refine((val) => val === "" || isMatch(val, "yyyy-MM-dd"), {
        message: "Start date is not in the correct format (YYYY-MM-DD)",
      }),
    endDate: z
      .string()
      .trim()
      .refine((val) => val === "" || isMatch(val, "yyyy-MM-dd"), {
        message: "End date is not in the correct format (YYYY-MM-DD)",
      }),

    tags: z.string(),
    creators: z.string(),
  });

  const { data } = $props();
  const apiClient = getApiClient();

  $effect(() => {
    reset({
      data: {
        type: data.media.type as MediaType,

        title: data.media.title,
        description: data.media.description ?? "",

        score: data.media.score ?? 0.0,
        status: data.media.status as MediaStatus,
        rating: data.media.rating as MediaRating,
        airingSeason: data.media.airingSeason ?? "",

        startDate: data.media.startDate ?? "",
        endDate: data.media.endDate ?? "",

        tags: data.media.tags.join(","),
        creators: data.media.creators.join(","),
      },
    });
  });

  const { form, errors, enhance, reset, submitting } = superForm(
    defaults(zod(Schema)),
    {
      SPA: true,
      validators: zod(Schema),
      dataType: "json",
      resetForm: true,
      async onUpdate({ form }) {
        if (form.valid) {
          const formData = form.data;
          const res = await apiClient.editMedia(data.media.id, {
            type: formData.type,

            title: formData.title,
            description: formData.description,

            score: formData.score,
            status: formData.status,
            rating: formData.rating,
            airingSeason: formData.airingSeason,

            startDate: formData.startDate,
            endDate: formData.endDate,

            tags: formData.tags
              .split(",")
              .map((s) => s.trim())
              .filter((s) => s !== ""),
            creators: formData.creators
              .split(",")
              .map((s) => s.trim())
              .filter((s) => s !== ""),
          });
          if (!res.success) {
            return handleApiError(res.error);
          }

          toast.success("Successfully updated media");
          invalidateAll();
        }
      },
    },
  );
</script>

<form class="flex flex-col gap-4 px-[1px]" use:enhance>
  <FormItem>
    <Label for="type">Type</Label>

    <Select.Root type="single" bind:value={$form.type} allowDeselect={false}>
      <Select.Trigger>
        {mediaTypes.find((f) => f.value === $form.type)?.label ??
          "Select type"}
      </Select.Trigger>
      <Select.Content>
        {#each mediaTypes as type (type.value)}
          <Select.Item value={type.value} label={type.label} />
        {/each}
      </Select.Content>
    </Select.Root>

    <Errors errors={$errors.type} />
  </FormItem>

  <FormItem>
    <Label for="title">Title</Label>
    <Input id="title" name="title" type="text" bind:value={$form.title} />
    <Errors errors={$errors.title} />
  </FormItem>

  <FormItem>
    <Label for="description">Description</Label>
    <textarea
      id="description"
      name="description"
      class="flex h-24 w-full resize-none rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
      rows={4}
      bind:value={$form.description}
    ></textarea>
    <Errors errors={$errors.description} />
  </FormItem>

  <FormItem>
    <Label for="score">Score</Label>
    <Input
      id="score"
      name="score"
      type="number"
      step={0.01}
      bind:value={$form.score}
    />
    <Errors errors={$errors.score} />
  </FormItem>

  <FormItem>
    <Label for="status">Status</Label>

    <Select.Root type="single" bind:value={$form.status} allowDeselect={false}>
      <Select.Trigger>
        {mediaStatus.find((f) => f.value === $form.status)?.label ??
          "Select status"}
      </Select.Trigger>
      <Select.Content>
        {#each mediaStatus as status (status.value)}
          <Select.Item value={status.value} label={status.label} />
        {/each}
      </Select.Content>
    </Select.Root>

    <Errors errors={$errors.status} />
  </FormItem>

  <FormItem>
    <Label for="rating">Rating</Label>

    <Select.Root type="single" bind:value={$form.rating} allowDeselect={false}>
      <Select.Trigger>
        {mediaRatings.find((f) => f.value === $form.rating)?.label ??
          "Select rating"}
      </Select.Trigger>
      <Select.Content>
        {#each mediaRatings as rating (rating.value)}
          <Select.Item value={rating.value} label={rating.label} />
        {/each}
      </Select.Content>
    </Select.Root>

    <Errors errors={$errors.rating} />
  </FormItem>

  <FormItem>
    <Label for="airingSeason">Airing Season</Label>
    <Input
      id="airingSeason"
      name="airingSeason"
      type="text"
      bind:value={$form.airingSeason}
    />
    <Errors errors={$errors.airingSeason} />
  </FormItem>

  <FormItem>
    <Label for="startDate">Start Date (YYYY-MM-DD)</Label>
    <Input
      id="startDate"
      name="startDate"
      type="text"
      bind:value={$form.startDate}
    />
    <Errors errors={$errors.startDate} />
  </FormItem>

  <FormItem>
    <Label for="endDate">End Date (YYYY-MM-DD)</Label>
    <Input
      id="endDate"
      name="endDate"
      type="text"
      bind:value={$form.endDate}
    />
    <Errors errors={$errors.endDate} />
  </FormItem>

  <FormItem>
    <Label for="tags">Tags (comma seperated)</Label>
    <Input id="tags" name="tags" type="text" bind:value={$form.tags} />
    <Errors errors={$errors.tags} />
  </FormItem>

  <FormItem>
    <Label for="creators">Creators (comma seperated)</Label>
    <Input
      id="creators"
      name="creators"
      type="text"
      bind:value={$form.creators}
    />
    <Errors errors={$errors.creators} />
  </FormItem>

  <Button type="submit" disabled={$submitting}>
    Update
    {#if $submitting}
      <Spinner />
    {/if}
  </Button>
</form>
