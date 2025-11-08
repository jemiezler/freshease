export type DialogProps = {
	open: boolean;
	onOpenChange: (open: boolean) => void;
	onSaved: () => Promise<void>;
};

export type EditDialogProps = {
	id: string;
	onOpenChange: (open: boolean) => void;
	onSaved: () => Promise<void>;
};

