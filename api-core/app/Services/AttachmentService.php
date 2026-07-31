<?php

declare(strict_types = 1);

namespace App\Services;

use App\Models\Attachment;
use Illuminate\Http\UploadedFile;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Queue;
use Illuminate\Support\Facades\Storage;
use Illuminate\Support\Str;

class AttachmentService
{
    public function __construct(
        private readonly UploadedFile $file,
        private readonly string $directory = '/attachments',
    ) {}

    private function getFilename(): string
    {
        return Str::uuid()->toString() . '.' . $this->file->getClientOriginalExtension();
    }

    private function getPath(): string
    {
        return "{$this->directory}/{$this->getFilename()}";
    }

    private function alreadyExists(): bool
    {
        return Storage::exists($this->getPath());
    }

    private function upload(): string | bool
    {
        return $this->file->storeAs($this->directory, $this->getFilename());
    }

    public function store(): void
    {
        if ($this->alreadyExists()) {
            throw new \InvalidArgumentException('File already exists');
        }

        try {
            $path = $this->upload();

            DB::transaction(function () use ($path) {
                $attachment = Attachment::query()->create([
                    'filename'  => $this->file->getClientOriginalName(),
                    'mime_type' => $this->file->getClientMimeType(),
                    'path'      => $path,
                    'user_id'   => auth()->id(),
                    'size'      => $this->file->getSize(),
                ]);

                DB::afterCommit(
                    fn() => Queue::connection('rabbitmq')
                        ->pushRaw(
                            json_encode(['id' => $attachment->id, 'path' => $path]),
                            'image-processor'
                        )
                );
            });
        } catch (\Exception $e) {
            if ($path) {
                Storage::delete($path);
            }
            report($e);

            throw new \RuntimeException('Failed to upload file: ' . $e->getMessage());
        }
    }
}
