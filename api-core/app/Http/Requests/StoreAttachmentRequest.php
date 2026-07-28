<?php

declare(strict_types = 1);

namespace App\Http\Requests;

use Illuminate\Foundation\Http\FormRequest;

class StoreAttachmentRequest extends FormRequest
{
    public function rules(): array
    {
        return [
            'file' => ['required', 'file', 'extensions:zip', 'mimes:application/zip', 'max:10240'],
        ];
    }
}
