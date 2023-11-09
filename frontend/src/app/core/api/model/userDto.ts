import {BingoCardDto} from './bingoCardDto';

export class UserDto {
    'ID': number;
    'github_id': string;
    'google_id': string;
    'reddit_id': string;
    'avatar_url': string;
    'name': string;
    'github_url': string;
    'reddit_url': string;
    'bingo_cards': BingoCardDto[];
}
