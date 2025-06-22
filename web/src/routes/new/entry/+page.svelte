<script lang="ts">
  import FormItem from "$lib/components/FormItem.svelte";
  import { Input, Label } from "@nanoteck137/nano-ui";
</script>

<p>New Entry</p>

<FormItem>
  <Label for="username">Username</Label>
  <Input id="username" name="username" type="text" />
  <!-- <Errors errors={$errors.username} /> -->
</FormItem>
