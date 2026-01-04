<script lang="ts">
	import Title from "$lib/components/Title.svelte";
	import Button from "$lib/components/Button.svelte";
	import SecondaryButton from "$lib/components/SecondaryButton.svelte";
	import TextInput from "$lib/components/TextInput.svelte";
	import PasswordInput from "$lib/components/PasswordInput.svelte";
	import Checkbox from "$lib/components/Checkbox.svelte";
	import Avatar from "$lib/components/Avatar.svelte";
	import LineSeparator from "$lib/components/LineSeparator.svelte";
	import AlertDialog from "$lib/components/AlertDialog.svelte";
	import BackNavBar from "$lib/components/BackNavBar.svelte";
	import Card from "$lib/components/Card.svelte";
	import Badge from "$lib/components/Badge.svelte";
	import Alert from "$lib/components/Alert.svelte";
	import Tag from "$lib/components/Tag.svelte";
	import Divider from "$lib/components/Divider.svelte";
	import Spinner from "$lib/components/Spinner.svelte";
	import Progress from "$lib/components/Progress.svelte";
	import Tooltip from "$lib/components/Tooltip.svelte";

	let dialogOpen = false;
	let checkboxChecked = false;
	let inputValue = "";
	let passwordValue = "";
	let isLoading = false;
	let progressValue = 65;
	let removableTags = ["React", "Svelte", "Vue"];

	function openDialog() {
		dialogOpen = true;
	}

	function closeDialog() {
		dialogOpen = false;
	}

	function handleConfirm() {
		isLoading = true;
		setTimeout(() => {
			isLoading = false;
			closeDialog();
		}, 2000);
	}

	function removeTag(index: number) {
		removableTags = removableTags.filter((_, i) => i !== index);
	}
</script>

<div class="min-h-screen bg-gradient-to-br from-blue-50 to-slate-200">
	<BackNavBar backPath="/">
		<span slot="first">Component Demo</span>
	</BackNavBar>

	<div class="max-w-2xl mx-auto p-8 space-y-8">
		<!-- Title Section -->
		<section>
			<Title
				title="Component Library"
				description="A complete showcase of all UI components with your custom color palette"
			/>
		</section>

		<LineSeparator />

		<!-- Buttons -->
		<section class="space-y-4">
			<h3 class="text-slate-900 text-lg font-semibold">Buttons</h3>
			<div class="flex gap-4 flex-wrap">
				<Button size="md" on:click={openDialog}>
					Open Dialog
				</Button>
				<Button size="sm">
					Small Button
				</Button>
				<SecondaryButton size="md">
					Secondary
				</SecondaryButton>
			</div>
		</section>

		<LineSeparator />

		<!-- Cards -->
		<section class="space-y-4">
			<h3 class="text-slate-900 text-lg font-semibold">Cards</h3>
			<Card>
				<h4 class="text-slate-900 font-semibold mb-2">Card Title</h4>
				<p class="text-slate-600">This is a content card that can be used to group related information together.</p>
			</Card>
		</section>

		<LineSeparator />

		<!-- Badges -->
		<section class="space-y-4">
			<h3 class="text-slate-900 text-lg font-semibold">Badges</h3>
			<div class="flex gap-3 flex-wrap">
				<Badge variant="primary">Primary</Badge>
				<Badge variant="secondary">Secondary</Badge>
				<Badge variant="danger">Danger</Badge>
				<Badge size="md" variant="primary">Medium Badge</Badge>
			</div>
		</section>

		<LineSeparator />

		<!-- Tags -->
		<section class="space-y-4">
			<h3 class="text-slate-900 text-lg font-semibold">Tags</h3>
			<div class="flex gap-2 flex-wrap">
				{#each removableTags as tag, i}
					<Tag removable onRemove={() => removeTag(i)}>
						{tag}
					</Tag>
				{/each}
			</div>
		</section>

		<LineSeparator />

		<!-- Alerts -->
		<section class="space-y-4">
			<h3 class="text-slate-900 text-lg font-semibold">Alerts</h3>
			<Alert type="info" title="Information">
				This is an informational alert message.
			</Alert>
			<Alert type="success" title="Success" dismissible>
				Your action was completed successfully!
			</Alert>
			<Alert type="warning" title="Warning">
				Please review this warning message carefully.
			</Alert>
			<Alert type="error" title="Error" dismissible>
				An error has occurred. Please try again.
			</Alert>
		</section>

		<LineSeparator />

		<!-- Progress -->
		<section class="space-y-4">
			<h3 class="text-slate-900 text-lg font-semibold">Progress Bar</h3>
			<Progress value={progressValue} />
		</section>

		<LineSeparator />

		<!-- Spinner & Tooltip -->
		<section class="space-y-4">
			<h3 class="text-slate-900 text-lg font-semibold">Spinner & Tooltip</h3>
			<div class="flex gap-6 items-center">
				<div class="flex flex-col gap-2 items-center">
					<Spinner size="sm" />
					<p class="text-xs text-slate-600">Small</p>
				</div>
				<div class="flex flex-col gap-2 items-center">
					<Spinner size="md" />
					<p class="text-xs text-slate-600">Medium</p>
				</div>
				<div class="flex flex-col gap-2 items-center">
					<Spinner size="lg" />
					<p class="text-xs text-slate-600">Large</p>
				</div>
				<Tooltip text="This is a helpful tooltip">
					<span class="px-3 py-2 bg-sky-400 text-white rounded cursor-help">
						Hover me
					</span>
				</Tooltip>
			</div>
		</section>

		<LineSeparator />

		<!-- Text Input -->
		<section class="space-y-4">
			<h3 class="text-slate-900 text-lg font-semibold">Text Input</h3>
			<TextInput
				placeholderText="Enter some text..."
				on:textchange={(e) => (inputValue = e.detail.value)}
			/>
			<p class="text-slate-600 text-sm">Value: {inputValue || "empty"}</p>
		</section>

		<LineSeparator />

		<!-- Password Input -->
		<section class="space-y-4">
			<h3 class="text-slate-900 text-lg font-semibold">Password Input</h3>
			<PasswordInput
				placeholderText="Enter your password..."
				bind:value={passwordValue}
				on:input={(e) => (passwordValue = e.detail.value)}
			/>
			<p class="text-slate-600 text-sm">Password: {passwordValue ? "●".repeat(passwordValue.length) : "empty"}</p>
		</section>

		<LineSeparator />

		<!-- Checkbox -->
		<section class="space-y-4">
			<h3 class="text-slate-900 text-lg font-semibold">Checkbox</h3>
			<div class="flex items-center gap-4">
				<Checkbox bind:checked={checkboxChecked} />
				<span class="text-slate-900">Checkbox is {checkboxChecked ? "checked" : "unchecked"}</span>
			</div>
		</section>

		<LineSeparator />

		<!-- Avatar -->
		<section class="space-y-4">
			<h3 class="text-slate-900 text-lg font-semibold">Avatar</h3>
			<div class="flex gap-4">
				<Avatar
					src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
					width="80"
					height="80"
				/>
				<Avatar
					src="https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?ixlib=rb-1.2.1&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
					width="60"
					height="60"
				/>
			</div>
		</section>

		<LineSeparator />

		<!-- Color Palette Reference -->
		<section class="space-y-4">
			<h3 class="text-slate-900 text-lg font-semibold">Color Palette</h3>
			<div class="grid grid-cols-5 gap-4">
				<div class="space-y-2">
					<div class="w-full h-20 rounded bg-blue-50 border border-slate-300"></div>
					<p class="text-xs text-center text-slate-600">E8EEF2</p>
				</div>
				<div class="space-y-2">
					<div class="w-full h-20 rounded bg-stone-300 border border-slate-300"></div>
					<p class="text-xs text-center text-slate-600">D6C9C9</p>
				</div>
				<div class="space-y-2">
					<div class="w-full h-20 rounded bg-slate-300 border border-slate-300"></div>
					<p class="text-xs text-center text-slate-600">C7D3DD</p>
				</div>
				<div class="space-y-2">
					<div class="w-full h-20 rounded bg-sky-400 border border-slate-300"></div>
					<p class="text-xs text-center text-slate-600">77B6EA</p>
				</div>
				<div class="space-y-2">
					<div class="w-full h-20 rounded bg-slate-900 border border-slate-300"></div>
					<p class="text-xs text-center text-slate-600">37393A</p>
				</div>
			</div>
		</section>
	</div>

	<!-- Alert Dialog -->
	<AlertDialog
		isOpen={dialogOpen}
		title="Confirm Action"
		description="This is a demo dialog. Click confirm to see the loading state."
		confirmText="Confirm"
		cancelText="Cancel"
		showAction={true}
		isDangerous={false}
		isLoading={isLoading}
		doConfirm={handleConfirm}
		doCancel={closeDialog}
	/>
</div>