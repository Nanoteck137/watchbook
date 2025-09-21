<script lang="ts">
  import { getApiClient, handleApiError } from "$lib";
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Dialog, Input, Label, Select } from "@nanoteck137/nano-ui";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";
  import Spinner from "$lib/components/Spinner.svelte";
  import type { MediaRelease } from "$lib/api/types";
  import toast from "svelte-5-french-toast";
  import { invalidateAll } from "$app/navigation";
  import { Trash } from "lucide-svelte";
  import ConfirmBox from "$lib/components/ConfirmBox.svelte";
  import {
    MediaReleaseTypeEnum,
    mediaReleaseTypes,
    type MediaReleaseType,
  } from "$lib/api-types";
  import { isValid, parseISO } from "date-fns";

  function isRFC3339(str: string): boolean {
    const date = parseISO(str);
    return isValid(date) && date.toISOString() === date.toISOString(); // ensures it round-trips
  }

  const Schema = z.object({
    type: MediaReleaseTypeEnum.default("not-confirmed"),
    startDate: z.string().refine((s) => isRFC3339(s)),
    numExpectedParts: z.number().min(0),
    intervalDays: z.number().min(0),
    delayDays: z.number().min(0),
  });

  export type Props = {
    open: boolean;
    mediaId: string;
    release?: MediaRelease;
  };

  let { open = $bindable(), mediaId, release }: Props = $props();
  const apiClient = getApiClient();

  let openDeleteConfirmModal = $state(false);

  $effect(() => {
    if (open) {
      reset({
        data: {
          type: (release?.releaseType as MediaReleaseType) ?? "not-confirmed",
          startDate: release?.startDate ?? "",
          numExpectedParts: release?.delayDays ?? 0,
          intervalDays: release?.delayDays ?? 0,
          delayDays: release?.delayDays ?? 0,
        },
      });
    }
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
          const data = form.data;

          const res = await apiClient.setMediaRelease(mediaId, {
            releaseType: data.type,
            startDate: data.startDate,
            numExpectedParts: data.numExpectedParts,
            intervalDays: data.intervalDays,
            delayDays: data.delayDays,
          });
          if (!res.success) {
            return handleApiError(res.error);
          }

          open = false;

          toast.success("Successfully updated list");
          invalidateAll();
        }
      },
    },
  );
</script>

<Dialog.Root bind:open>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Set release</Dialog.Title>
    </Dialog.Header>

    <form class="flex flex-col gap-4" use:enhance>
      <FormItem>
        <Label for="type">Type</Label>

        <div class="flex items-center gap-2">
          <Select.Root
            type="single"
            bind:value={$form.type}
            allowDeselect={false}
          >
            <Select.Trigger>
              {mediaReleaseTypes.find((f) => f.value === $form.type)?.label ??
                "Select type"}
            </Select.Trigger>
            <Select.Content>
              {#each mediaReleaseTypes as ty (ty.value)}
                <Select.Item value={ty.value} label={ty.label} />
              {/each}
            </Select.Content>
          </Select.Root>

          <Button
            variant="destructive"
            onclick={() => {
              openDeleteConfirmModal = true;
            }}
          >
            <Trash />
          </Button>
        </div>

        <Errors errors={$errors.type} />
      </FormItem>

      <FormItem>
        <Label for="startDate">
          Start Date (RFC3339) (2025-07-10T13:00:00Z)
        </Label>
        <Input
          id="startDate"
          name="startDate"
          type="text"
          bind:value={$form.startDate}
        />
        <Errors errors={$errors.startDate} />
      </FormItem>

      <FormItem>
        <Label for="numExpectedParts">Expected parts</Label>
        <Input
          id="numExpectedParts"
          name="numExpectedParts"
          type="number"
          min={0}
          bind:value={$form.numExpectedParts}
        />
        <Errors errors={$errors.numExpectedParts} />
      </FormItem>

      <FormItem>
        <Label for="intervalDays">Interval (days)</Label>
        <Input
          id="intervalDays"
          name="intervalDays"
          type="number"
          min={0}
          bind:value={$form.intervalDays}
        />
        <Errors errors={$errors.intervalDays} />
      </FormItem>

      <FormItem>
        <Label for="delayDays">Delay (days)</Label>
        <Input
          id="delayDays"
          name="delayDays"
          type="number"
          min={0}
          bind:value={$form.delayDays}
        />
        <Errors errors={$errors.delayDays} />
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
          Set Release
          {#if $submitting}
            <Spinner />
          {/if}
        </Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>

<ConfirmBox
  bind:open={openDeleteConfirmModal}
  title="Remove release?"
  description="Are you sure you want to remove this release? This action cannot be undone."
  onResult={async () => {
    const res = await apiClient.deleteMediaRelease(mediaId);
    if (!res.success) {
      return handleApiError(res.error);
    }

    open = false;
    invalidateAll();
  }}
/>
