<script lang="ts">
  import { getApiClient, handleApiError } from "$lib";
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Input, Label } from "@nanoteck137/nano-ui";
  import SuperDebug from "sveltekit-superforms";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";
  import toast from "svelte-5-french-toast";
  import { goto } from "$app/navigation";

  const apiClient = getApiClient();

  let loading = $state(false);

  const formSchema = z.object({
    title: z.string().trim().min(1),
    malIds: z.string().trim(),
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
          const res = await apiClient.createCollection({
            collectionType: "anime",
            name: f.data.title,
          });
          if (!res.success) {
            return handleApiError(res.error);
          }

          const collectionId = res.data.id;

          const splits = f.data.malIds.split(",");
          for (let id of splits) {
            id = id.trim();
            if (id.length === 0) {
              return;
            }

            const animeDataRes =
              await apiClient.providerMyAnimeListGetAnime(id);
            if (!animeDataRes.success) {
              // TODO(patrik): Better handling of error
              handleApiError(animeDataRes.error);
              continue;
            }

            const animeData = animeDataRes.data;

            const title = animeData.titleEnglish ?? animeData.title;
            const mediaRes = await apiClient.createMedia({
              title: title,
              description: animeData.description,
              mediaType: animeData.mediaType,
              score: animeData.score ?? 0,
              status: animeData.status,
              rating: animeData.rating,
              airingSeason: animeData.airingSeason,
              startDate: animeData.startDate ?? "",
              endDate: animeData.endDate ?? "",
              studios: animeData.studios,
              tags: animeData.tags,
              tmdbId: "",
              malId: animeData.malId,
              anilistId: "",
              partCount: animeData.episodeCount ?? 0,
              coverUrl: animeData.coverImageUrl,
              bannerUrl: "",
              logoUrl: "",
              collectionId: collectionId,
              collectionName: title,
            });
            if (!mediaRes.success) {
              return handleApiError(mediaRes.error);
            }
          }

          toast.success("Successfully created media item");
          goto(`/collections/${collectionId}`, { invalidateAll: true });
        }

        // if (f.valid) {
        //   const res = await apiClient.createMedia({
        //     title: f.data.title,
        //     description: f.data.description,
        //     mediaType: f.data.type,
        //     score: f.data.score,
        //     status: f.data.status,
        //     rating: f.data.rating,
        //     airingSeason: f.data.airingSeason,
        //     startDate: f.data.startDate,
        //     endDate: f.data.endDate,
        //     studios: f.data.studios.split(","),
        //     tags: f.data.tags.split(","),
        //     tmdbId: f.data.tmdbId,
        //     malId: f.data.malId,
        //     anilistId: f.data.anilistId,
        //     partCount: f.data.generateParts ? f.data.partCount : 0,
        //     coverUrl: f.data.coverUrl,
        //     bannerUrl: "",
        //     logoUrl: "",
        //   });
        //   if (!res.success) {
        //     return handleApiError(res.error);
        //   }

        // TODO(patrik): Navigate to media item
        // toast.success("Successfully created media item");
        // goto(`/media/${res.data.id}`, { invalidateAll: true });
        // }
      },
    },
  );

  /*$effect(() => {
    $form.type = parseMediaType(data?.mediaType);

    $form.tmdbId = "";
    $form.malId = data?.malId ?? "";
    $form.anilistId = "";

    $form.title = pickTitle(data ?? { title: "", titleEnglish: "" });
    $form.description = data?.description ?? "";

    $form.score = data?.score ?? 0.0;
    $form.status = parseMediaStatus(data?.status);
    $form.rating = parseMediaRating(data?.rating);
    $form.airingSeason = data?.airingSeason ?? "";

    $form.startDate = data?.startDate ?? "";
    $form.endDate = data?.endDate ?? "";

    $form.tags = data?.tags.join(",") ?? "";
    $form.studios = data?.studios.join(",") ?? "";

    $form.coverUrl = data?.coverImageUrl ?? "";

    $form.partCount = data?.episodeCount ?? 0;
    $form.generateParts = true;
  });*/
</script>

<p>Loading: {loading}</p>

<form method="POST" class="w-2/3 space-y-6" use:enhance>
  <FormItem>
    <Label for="title" aria-invalid={$errors.title ? "true" : undefined}>
      Serie Title
    </Label>
    <Input id="title" name="title" type="text" bind:value={$form.title} />
    {#if $errors.title}
      <span class="text-red-200">{$errors.title}</span>
    {/if}
  </FormItem>

  <FormItem>
    <Label for="malIds" aria-invalid={$errors.malIds ? "true" : undefined}>
      MAL Ids
    </Label>
    <Input id="malIds" name="malIds" type="text" bind:value={$form.malIds} />
    {#if $errors.malIds}
      <span class="text-red-200">{$errors.malIds}</span>
    {/if}
  </FormItem>

  <Button type="submit" disabled={$submitting}>Submit</Button>
</form>

<SuperDebug data={$form} />
