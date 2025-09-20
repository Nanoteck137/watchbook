<script lang="ts">
  import { getApiClient, handleApiError } from "$lib";
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Dialog, Input, Label, Select } from "@nanoteck137/nano-ui";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";
  import { isMatch } from "date-fns";
  import Spinner from "$lib/components/Spinner.svelte";
  import type { Media } from "$lib/api/types";
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
  } from "$lib/api-types";

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
  type SchemaTy = z.infer<typeof Schema>;

  export type Props = {
    open: boolean;
    media: Media;
  };

  let { open = $bindable(), media }: Props = $props();
  const apiClient = getApiClient();

  $effect(() => {
    if (open) {
      reset({
        data: {
          type: media.mediaType as MediaType,

          title: media.title,
          description: media.description ?? "",

          score: media.score ?? 0.0,
          status: media.status as MediaStatus,
          rating: media.rating as MediaRating,
          airingSeason: media.airingSeason ?? "",

          startDate: media.startDate ?? "",
          endDate: media.endDate ?? "",

          tags: media.tags.join(","),
          creators: media.creators.join(","),
        },
      });
    }
  });

  async function submit(data: SchemaTy) {
    console.log(data);

    const res = await apiClient.editMedia(media.id, {
      mediaType: data.type,

      title: data.title,
      description: data.description,

      score: data.score,
      status: data.status,
      rating: data.rating,
      airingSeason: data.airingSeason,

      startDate: data.startDate,
      endDate: data.endDate,

      tags: data.tags
        .split(",")
        .map((s) => s.trim())
        .filter((s) => s !== ""),
      creators: data.creators
        .split(",")
        .map((s) => s.trim())
        .filter((s) => s !== ""),
    });
    if (!res.success) {
      return handleApiError(res.error);
    }

    open = false;

    toast.success("Successfully updated media");
    invalidateAll();
  }

  const { form, errors, enhance, reset, submitting } = superForm(
    defaults(zod(Schema)),
    {
      SPA: true,
      validators: zod(Schema),
      dataType: "json",
      resetForm: true,
      async onUpdate({ form }) {
        if (form.valid) {
          await submit(form.data);
        }
      },
    },
  );
</script>

<Dialog.Root bind:open>
  <Dialog.Content class="max-h-[420px] overflow-y-scroll">
    <Dialog.Header>
      <Dialog.Title>Create new media</Dialog.Title>
    </Dialog.Header>

    <form class="flex flex-col gap-4 px-[1px]" use:enhance>
      <FormItem>
        <Label for="type">Type</Label>

        <Select.Root
          type="single"
          bind:value={$form.type}
          allowDeselect={false}
        >
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

        <Select.Root
          type="single"
          bind:value={$form.status}
          allowDeselect={false}
        >
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

        <Select.Root
          type="single"
          bind:value={$form.rating}
          allowDeselect={false}
        >
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

      <Dialog.Footer class="gap-2 sm:gap-0">
        <Button
          variant="outline"
          onclick={() => {
            open = false;
            reset();
          }}
        >
          Close
        </Button>

        <Button type="submit" disabled={$submitting}>
          Update
          {#if $submitting}
            <Spinner />
          {/if}
        </Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
