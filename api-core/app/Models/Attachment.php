<?php

declare(strict_types = 1);

namespace App\Models;

use App\AttachmentUploadStatus;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

class Attachment extends Model
{
    protected $fillable = [
        'user_id',
        'filename',
        'path',
        'mime_type',
        'size',
        'status',
    ];

    protected $casts = [
        'status' => AttachmentUploadStatus::class,
    ];

    public function user(): BelongsTo
    {
        return $this->belongsTo(User::class);
    }
}
